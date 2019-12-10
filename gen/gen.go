package gen

import (
	"fmt"
	"io"
	"strings"

	"github.com/hfaulds/tracer/parse/types"
)

func GenStruct(b io.Writer, importMap ImportMap, strct types.Struct) {
	fmt.Fprintf(b, "\n\ntype %s struct {\n", strct.Name)
	for _, attr := range strct.Attrs {
		fmt.Fprintf(b, "\t%s %s\n", resolveParam(importMap, attr))
	}
	fmt.Fprintf(b, "}\n")
}

func GenMethod(b io.Writer, importMap ImportMap, strct *types.Struct, method types.Method, callback func(b io.Writer)) {
	fmt.Fprint(b, "\n\nfunc ")
	if strct != nil {
		fmt.Fprintf(b, "(t %s) ", strct.Name)
	}
	generateMethodSig(b, "", method.Name, resolveParams(importMap, method.Params), resolveParams(importMap, method.Returns))
	fmt.Fprint(b, " {\n")
	callback(b)
	fmt.Fprintf(b, "\n}\n")
}

func resolveParams(importMap map[string]string, params []types.Param) []string {
	resolved := make([]string, 0, len(params))
	for _, p := range params {
		resolved = append(resolved, resolveParam(importMap, p))
	}
	return resolved
}

func resolveParam(importMap map[string]string, p types.Param) string {
	switch tp := p.(type) {
	case types.BasicParam:
		return tp.Typ
	case types.NamedParam:
		if tp.Pkg != "" {
			if alias, ok := importMap[tp.Pkg]; ok {
				return fmt.Sprintf("%s.%s", alias, tp.Typ)
			} else {
				return tp.Typ
			}
		}
		return tp.Typ
	case types.ArrayParam:
		return fmt.Sprintf("[%d]%s", tp.Length, resolveParam(importMap, tp.Typ))
	case types.SliceParam:
		return fmt.Sprintf("[]%s", resolveParam(importMap, tp.Typ))
	case types.PointerParam:
		return fmt.Sprintf("*%s", resolveParam(importMap, tp.Typ))
	case types.MapParam:
		return fmt.Sprintf("map[%s]%s", resolveParam(importMap, tp.Key), resolveParam(importMap, tp.Elem))
	case types.InterfaceParam:
		var b strings.Builder
		if len(tp.Methods) == 0 {
			fmt.Fprint(&b, "interface{}")
		} else if len(tp.Methods) == 1 {
			fmt.Fprint(&b, "interface{ ")
			m := tp.Methods[0]
			params := resolveParams(importMap, m.Params)
			returns := resolveParams(importMap, m.Returns)
			generateMethodSig(&b, "", m.Name, params, returns)
			fmt.Fprint(&b, " }")
		} else {
			fmt.Fprint(&b, "interface {")
			for _, m := range tp.Methods {
				fmt.Fprint(&b, "\n\t")
				params := resolveParams(importMap, m.Params)
				returns := resolveParams(importMap, m.Returns)
				generateMethodSig(&b, "", m.Name, params, returns)
			}
			fmt.Fprint(&b, "\n},\n")
		}
		return b.String()
	default:
		return "<unsupported>"
	}
}

func generateMethodSig(b io.Writer, implementor, methodName string, params, returns []string) {
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
