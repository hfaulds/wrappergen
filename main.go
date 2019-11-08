package main

import (
	"errors"
	"flag"
	"fmt"
	"go/types"
	"reflect"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/packages"
)

func main() {
	flag.Parse()
	args := flag.Args()
	p, err := getPackage(args[0])
	if err != nil {
		panic(err)
	}
	fmt.Println("types", spew.Sdump(getTypes(p)))
}

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

func getTypes(p *packages.Package) []Interface {
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
	return interfaces
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
		fmt.Println(reflect.TypeOf(v.Type()))
		p := param{
			Name: v.Name(),
			Type: types.TypeString(v.Type(), nil),
		}
		params = append(params, p)
	}
	return params
}
