package types

type Package struct {
	Name       string
	PkgPath    string
	Interfaces []Interface
	Structs    []Struct
}

func (pkg *Package) FindInterface(name string) (Interface, bool) {
	for _, iface := range pkg.Interfaces {
		if iface.Name == name {
			return iface, true
		}
	}
	return Interface{}, false
}

func (pkg *Package) FindStruct(name string) (Struct, bool) {
	for _, strct := range pkg.Structs {
		if strct.Name == name {
			return strct, true
		}
	}
	return Struct{}, false
}

type Interface struct {
	Name    string
	File    string
	Methods []Method
}

type Struct struct {
	Name  string
	Attrs []Var
}

type Method struct {
	Name    string
	Params  []Param
	Returns []Param
}

type Var struct {
	Name string
	Type Param
}

type Param interface {
	Param()
}

type UnsupportedParam struct {
}

type BasicParam struct {
	Typ string
}

type NamedParam struct {
	Typ string
	Pkg string
}

type ArrayParam struct {
	Typ    Param
	Length int64
}

type SliceParam struct {
	Typ Param
}

type PointerParam struct {
	Typ Param
}

type MapParam struct {
	Key  Param
	Elem Param
}

type InterfaceParam struct {
	Methods []Method
}

func (p UnsupportedParam) Param() {}
func (p BasicParam) Param()       {}
func (p NamedParam) Param()       {}
func (p ArrayParam) Param()       {}
func (p SliceParam) Param()       {}
func (p PointerParam) Param()     {}
func (p MapParam) Param()         {}
func (p InterfaceParam) Param()   {}
