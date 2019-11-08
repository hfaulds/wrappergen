package main

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

func ParseDir(dir string) (*Package, error) {
	p, err := getPackage(dir)
	if err != nil {
		return nil, err
	}
	scope := p.Types.Scope()
	names := scope.Names()
	pkg := &Package{
		Name: p.Name,
	}
	for _, name := range names {
		obj := scope.Lookup(name)
		itype, ok := obj.Type().Underlying().(*types.Interface)
		if ok {
			file := p.Fset.File(obj.Pos())
			iface := Interface{
				Name:    name,
				File:    file.Name(),
				Methods: getMethods(itype),
			}
			pkg.Interfaces = append(pkg.Interfaces, iface)
		}
	}
	return pkg, nil
}

func getPackage(src string) (*packages.Package, error) {
	conf := packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedImports,
		Dir:  src,
	}
	pkgs, err := packages.Load(&conf)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, errors.New("No packages found")
	}
	if len(pkgs) > 1 {
		return nil, errors.New("More than one package was found")
	}
	return pkgs[0], nil
}

func getMethods(itype *types.Interface) []method {
	methods := make([]method, 0, itype.NumMethods())
	for i := 0; i < itype.NumMethods(); i++ {
		funcType := itype.Method(i)
		sig := funcType.Type().(*types.Signature)
		m := method{
			Name:    funcType.Name(),
			Params:  getParams(sig.Params()),
			Returns: getParams(sig.Results()),
		}
		methods = append(methods, m)
	}
	return methods
}

func getParams(tuple *types.Tuple) []param {
	params := make([]param, 0, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)

		params = append(params, getParam(v.Type()))
	}
	return params
}

func getParam(typ types.Type) param {
	switch t := typ.(type) {
	case *types.Basic:
		return basicParam{typ: t.Name()}
	case *types.Named:
		if obj := t.Obj(); obj != nil {
			var pkg string
			if p := obj.Pkg(); p != nil {
				pkg = p.Path()
			}
			return namedParam{typ: obj.Name(), pkg: pkg}
		}
	case *types.Array:
		return arrayParam{typ: getParam(t.Elem()), length: t.Len()}
	case *types.Slice:
		return sliceParam{typ: getParam(t.Elem())}
	case *types.Pointer:
		return pointerParam{typ: getParam(t.Elem())}
	case *types.Map:
		return mapParam{key: getParam(t.Key()), elem: getParam(t.Elem())}
	}
	return unsupportedParam{}
}
