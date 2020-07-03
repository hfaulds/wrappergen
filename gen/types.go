package gen

import (
	"github.com/hfaulds/tracer/parse/types"
)

type File struct {
	Structs []Struct
	Methods []Method
}

type Method struct {
	types.Method
	Callback func(b Builder, m types.Method)
}

type Struct struct {
	types.Struct
	Methods []Method
}
