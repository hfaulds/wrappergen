package gen

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/hfaulds/tracer/parse/types"
)

type ImportMap map[string]string

func GenerateImports(b io.Writer, importMap ImportMap) {
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

func BuildImportMap(pkg *types.Package) ImportMap {
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

func resolveMethodPackages(methods []types.Method) []string {
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

func resolvePackages(p types.Param) []string {
	switch tp := p.(type) {
	case types.NamedParam:
		if tp.Pkg == "" {
			return []string{}
		}
		return []string{tp.Pkg}
	case types.ArrayParam:
		return resolvePackages(tp.Typ)
	case types.SliceParam:
		return resolvePackages(tp.Typ)
	case types.PointerParam:
		return resolvePackages(tp.Typ)
	case types.MapParam:
		return append(resolvePackages(tp.Key), resolvePackages(tp.Elem)...)
	case types.InterfaceParam:
		return resolveMethodPackages(tp.Methods)
	default:
		return []string{}
	}
}
