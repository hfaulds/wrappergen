package main

import (
	"fmt"
	"sort"
	"strings"
)

func Generate(pkg *Package) string {
	var b strings.Builder

	/*
		package main
	*/
	fmt.Fprintf(&b, "package %s\n\n", pkg.Name)

	/*
		import (
			i0 "context"
			i1 "some/dependency"
		)
	*/
	importMap := buildImportMap(pkg)
	generateImports(&b, importMap)

	childSpanType := fmt.Sprintf("func(%s.Context) (%s.Context, interface{ Close() })", importMap["context"], importMap["context"])

	for _, iface := range pkg.Interfaces {
		// Skip interfaces where no methods have context.Context as an argument
		if shouldSkipInterface(iface) {
			continue
		}
		fmt.Fprint(&b, "\n")

		/*
			type traceExample struct {
				wrapped Example
			}
		*/
		structName := fmt.Sprintf("trace%s", iface.Name)
		fmt.Fprintf(&b, "type %s struct {\n", structName)
		fmt.Fprintf(&b, "\twrapped   %s\n", iface.Name)
		fmt.Fprintf(&b, "\tchildSpan %s\n", childSpanType)
		fmt.Fprintf(&b, "}")

		/* func NewExampleTracer(p0 Example) Example {
			return traceExample {
				wrapped p0,
			}
		}*/
		fmt.Fprint(&b, "\n\nfunc ")
		generateMethodSig(&b, "", fmt.Sprintf("New%sTracer", strings.Title(iface.Name)), []string{iface.Name, childSpanType}, []string{iface.Name})
		fmt.Fprint(&b, " {\n")
		fmt.Fprintf(&b, "\treturn %s{\n", structName)
		fmt.Fprintf(&b, "\t\twrapped:   p0,\n")
		fmt.Fprintf(&b, "\t\tchildSpan: p1,\n")
		fmt.Fprintf(&b, "\t}\n")
		fmt.Fprintf(&b, "}")

		/*
			func (t traceExample) Foo(p0 context.Context, p1) i1.example {
				ctx, span := t.childSpan(p0)
				defer span.Close()
				return t.wrapped.Foo(p0,p1)
			}
		*/
		for _, m := range iface.Methods {
			params := resolveParams(importMap, m.Params)
			returns := resolveParams(importMap, m.Returns)
			fmt.Fprint(&b, "\n\nfunc ")
			generateMethodSig(&b, structName, m.Name, params, returns)
			fmt.Fprint(&b, " {\n")
			// only add tracing if there a context
			offset, ok := getFirstContextParamOffset(m)
			if ok {
				fmt.Fprintf(&b, "\tctx, span := t.childSpan(p%d)\n", offset)
				fmt.Fprint(&b, "\tdefer span.Close()\n")
			}
			generateWrappedCall(&b, m, len(params), offset)
			fmt.Fprintf(&b, "}")
		}
		fmt.Fprint(&b, "\n")
	}
	return b.String()
}

func generateMethodSig(b *strings.Builder, implementor, methodName string, params, returns []string) {
	if implementor != "" {
		fmt.Fprintf(b, "(t %s) ", implementor)
	}
	fmt.Fprintf(b, "%s(", methodName)
	for i, param := range params {
		fmt.Fprintf(b, "p%d %s", i, param)
		if i < len(params)-1 {
			fmt.Fprint(b, ", ")
		}
	}
	fmt.Fprint(b, ")")
	if len(returns) > 0 {
		fmt.Fprint(b, " ")
	}
	if len(returns) > 1 {
		fmt.Fprint(b, "(")
	}
	for i, r := range returns {
		if i > 0 {
			fmt.Fprint(b, ", ")
		}
		fmt.Fprint(b, r)
	}
	if len(returns) > 1 {
		fmt.Fprint(b, ")")
	}
}

func buildImportMap(pkg *Package) map[string]string {
	importMap := map[string]string{}
	for _, i := range pkg.Interfaces {
		for _, p := range resolveMethodPackages(i.Methods) {
			if p == pkg.PkgPath {
				continue
			}
			if _, ok := importMap[p]; !ok {
				importMap[p] = fmt.Sprintf("i%d", len(importMap))
			}
		}
	}
	return importMap
}

func resolveMethodPackages(methods []method) []string {
	var pkgs []string
	for _, m := range methods {
		for _, p := range m.Params {
			pkgs = append(pkgs, resolvePackages(p)...)
		}
		for _, p := range m.Returns {
			pkgs = append(pkgs, resolvePackages(p)...)
		}
	}
	return pkgs
}

func resolvePackages(p param) []string {
	switch tp := p.(type) {
	case namedParam:
		return []string{tp.pkg}
	case arrayParam:
		return resolvePackages(tp.typ)
	case sliceParam:
		return resolvePackages(tp.typ)
	case pointerParam:
		return resolvePackages(tp.typ)
	case mapParam:
		return append(resolvePackages(tp.key), resolvePackages(tp.elem)...)
	case interfaceParam:
		return resolveMethodPackages(tp.methods)
	default:
		return []string{}
	}
}

func generateImports(b *strings.Builder, importMap map[string]string) {
	var imports []string
	for imp, alias := range importMap {
		imports = append(imports, fmt.Sprintf("import %s \"%s\"", alias, imp))
	}
	sort.Strings(imports)
	fmt.Fprintf(b, strings.Join(imports, "\n"))
	if len(imports) > 0 {
		fmt.Fprintf(b, "\n")
	}
}

func shouldSkipInterface(i Interface) bool {
	for _, m := range i.Methods {
		if _, ok := getFirstContextParamOffset(m); ok {
			return false
		}
	}
	return true
}

var contextNamedParam = namedParam{pkg: "context", typ: "Context"}

func getFirstContextParamOffset(m method) (int, bool) {
	for i, p := range m.Params {
		if np, ok := p.(namedParam); ok {
			if np == contextNamedParam {
				return i, true
			}
		}
	}
	return -1, false
}

func resolveParams(importMap map[string]string, params []param) []string {
	resolved := make([]string, 0, len(params))
	for _, p := range params {
		resolved = append(resolved, resolveParam(importMap, p))
	}
	return resolved
}

func resolveParam(importMap map[string]string, p param) string {
	switch tp := p.(type) {
	case basicParam:
		return tp.typ
	case namedParam:
		if tp.pkg != "" {
			if alias, ok := importMap[tp.pkg]; ok {
				return fmt.Sprintf("%s.%s", alias, tp.typ)
			} else {
				return tp.typ
			}
		}
		return tp.typ
	case arrayParam:
		return fmt.Sprintf("[%d]%s", tp.length, resolveParam(importMap, tp.typ))
	case sliceParam:
		return fmt.Sprintf("[]%s", resolveParam(importMap, tp.typ))
	case pointerParam:
		return fmt.Sprintf("*%s", resolveParam(importMap, tp.typ))
	case mapParam:
		return fmt.Sprintf("map[%s]%s", resolveParam(importMap, tp.key), resolveParam(importMap, tp.elem))
	case interfaceParam:
		var b strings.Builder
		fmt.Fprint(&b, "interface {")
		for _, m := range tp.methods {
			params := resolveParams(importMap, m.Params)
			returns := resolveParams(importMap, m.Returns)
			fmt.Fprint(&b, "\n\t")
			generateMethodSig(&b, "", m.Name, params, returns)
		}
		fmt.Fprint(&b, "\n},\n")
		return b.String()
	default:
		return "<unsupported>"
	}
}

func generateWrappedCall(b *strings.Builder, m method, numParams, offset int) {
	fmt.Fprint(b, "\t")
	if len(m.Returns) > 0 {
		fmt.Fprint(b, "return ")
	}
	fmt.Fprintf(b, "t.wrapped.%s(", m.Name)
	for i := 0; i < numParams; i++ {
		if i == offset {
			fmt.Fprint(b, "ctx")
		} else {
			fmt.Fprintf(b, "p%d", i)
		}
		if i != numParams-1 {
			fmt.Fprint(b, ", ")
		}
	}
	fmt.Fprint(b, ")\n")
}
