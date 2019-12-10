package parse

import (
	"errors"
	"fmt"
	gtypes "go/types"

	"golang.org/x/tools/go/packages"

	"github.com/hfaulds/tracer/parse/types"
)

func ParseDir(dir string) (*types.Package, error) {
	p, err := getPackage(dir)
	if err != nil {
		return nil, err
	}
	scope := p.Types.Scope()
	names := scope.Names()
	pkg := &types.Package{
		Name:    p.Name,
		PkgPath: p.PkgPath,
	}
	for _, name := range names {
		obj := scope.Lookup(name)
		itype, ok := obj.Type().Underlying().(*gtypes.Interface)
		if ok {
			file := p.Fset.File(obj.Pos())
			iface := types.Interface{
				Name:    name,
				File:    file.Name(),
				Methods: getMethods(itype),
			}
			pkg.Interfaces = append(pkg.Interfaces, iface)
		}
		stype, ok := obj.Type().Underlying().(*gtypes.Struct)
		if ok {
			strct := types.Struct{
				Name:  name,
				Attrs: getAttrs(stype),
			}
			pkg.Structs = append(pkg.Structs, strct)
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

func getMethods(itype *gtypes.Interface) []types.Method {
	methods := make([]types.Method, 0, itype.NumMethods())
	for i := 0; i < itype.NumMethods(); i++ {
		funcType := itype.Method(i)
		sig := funcType.Type().(*gtypes.Signature)
		m := types.Method{
			Name:    funcType.Name(),
			Params:  getParams(sig.Params()),
			Returns: getParams(sig.Results()),
		}
		methods = append(methods, m)
	}
	return methods
}

func getAttrs(stype *gtypes.Struct) map[string]types.Param {
	attrs := make(map[string]types.Param, stype.NumFields())
	for i := 0; i < stype.NumFields(); i++ {
		field := stype.Field(i)
		attrs[field.Name()] = getParam(field.Type())
	}
	return attrs
}

func getParams(tuple *gtypes.Tuple) []types.Param {
	params := make([]types.Param, 0, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)

		params = append(params, getParam(v.Type()))
	}
	return params
}

func getParam(typ gtypes.Type) types.Param {
	switch t := typ.(type) {
	case *gtypes.Basic:
		return types.BasicParam{Typ: t.Name()}
	case *gtypes.Named:
		if obj := t.Obj(); obj != nil {
			if obj.Name() == "returnType" {
				fmt.Println(obj.Name())
			}

			var pkg string
			if p := obj.Pkg(); p != nil {
				pkg = p.Path()
			}
			return types.NamedParam{Typ: obj.Name(), Pkg: pkg}
		}
	case *gtypes.Array:
		return types.ArrayParam{Typ: getParam(t.Elem()), Length: t.Len()}
	case *gtypes.Slice:
		return types.SliceParam{Typ: getParam(t.Elem())}
	case *gtypes.Pointer:
		return types.PointerParam{Typ: getParam(t.Elem())}
	case *gtypes.Map:
		return types.MapParam{Key: getParam(t.Key()), Elem: getParam(t.Elem())}
	case *gtypes.Interface:
		return types.InterfaceParam{Methods: getMethods(t)}
	}
	return types.UnsupportedParam{}
}
