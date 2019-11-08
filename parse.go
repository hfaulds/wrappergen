package main

import (
	"errors"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Interface struct {
	Name    string
	File    string
	Methods []method
}

type method struct {
	Name    string
	Params  []param
	Returns []param
}

type param struct {
	Name string
	Type string
}

func parseDir(dir string) ([]Interface, error) {
	p, err := getPackage(dir)
	if err != nil {
		return nil, err
	}
	scope := p.Types.Scope()
	names := scope.Names()
	var interfaces []Interface
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
			interfaces = append(interfaces, iface)
		}
	}
	return interfaces, nil
}

func getPackage(src string) (*packages.Package, error) {
	conf := packages.Config{
		Mode: packages.NeedTypes | packages.NeedImports,
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
			Name:   funcType.Name(),
			Params: getParams(sig),
		}
		methods = append(methods, m)
	}
	return methods
}

func getParams(sig *types.Signature) []param {
	tuple := sig.Params()
	params := make([]param, 0, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		p := param{
			Name: v.Name(),
			Type: types.TypeString(v.Type(), nil),
		}
		params = append(params, p)
	}
	return params
}
