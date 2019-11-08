package main

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

type param interface {
	Param()
}

type unsupportedParam struct {
}

type basicParam struct {
	typ string
}

type namedParam struct {
	typ string
	pkg string
}

type arrayParam struct {
	typ param
}

type sliceParam struct {
	typ param
}

type pointerParam struct {
	typ param
}

type mapParam struct {
	key  param
	elem param
}

type interfaceParam struct {
	methods []method
}

func (p unsupportedParam) Param() {}
func (p basicParam) Param()       {}
func (p namedParam) Param()       {}
func (p arrayParam) Param()       {}
func (p sliceParam) Param()       {}
func (p pointerParam) Param()     {}
func (p mapParam) Param()         {}
func (p interfaceParam) Param()   {}
