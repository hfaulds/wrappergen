package timing

import (
	"fmt"
	"strings"

	"github.com/hfaulds/wrappergen/gen"
	"github.com/hfaulds/wrappergen/parse/types"
)

var timingType = types.NamedParam{
	Pkg: "github.com/hfaulds/wrappergen/timing",
	Typ: "Timing",
}

var TimingWrapper = func(iface types.Interface) gen.Wrapper {
	return gen.Wrapper{
		Constructor: constructorName(iface),
		Arguments:   []types.NamedParam{timingType},
	}
}

func Gen(b gen.Builder, iface types.Interface) {
	timingStruct := gen.Struct{
		Struct: types.Struct {
			Name: fmt.Sprintf("time%s", iface.Name),
			Attrs: []types.Var{
				{Name: "wrapped", Type: types.NamedParam{Typ: iface.Name}},
				{Name: "timing", Type: timingType},
			},
		},
		Methods: make([]gen.Method, len(iface.Methods)),
	}

	for i, m := range iface.Methods {
		timingStruct.Methods[i] = gen.Method{
			Method: m,
			Callback: func(b gen.Builder, m types.Method) {
				offset, ok := getFirstContextParamOffset(m)
				if ok {
					b.WriteLine("timer := t.timing.Timer()")
					b.WriteLine("defer timer.End(p%d, \"%s\")", offset, m.Name)
				}
				numReturns := len(m.Returns)
				if numReturns > 0 {
					b.Write("return ")
				}
				b.Write("t.wrapped.%s(", m.Name)
				for i := 0; i < len(m.Params); i++ {
					b.Write("p%d", i)
					if i != len(m.Params)-1 {
						b.Write(", ")
					}
				}
				b.Write(")")
			},
		}
	}

	file := gen.File{
		Structs: []gen.Struct{ timingStruct },
		Methods: []gen.Method{
			{
				Method: types.Method{
					Name: constructorName(iface),
					Params: []types.Param{
						types.NamedParam{Typ: iface.Name},
						timingType,
					},
					Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
				},
				Callback: func(b gen.Builder, _ types.Method) {
					b.WriteLine("return %s{", timingStruct.Name)
					b.WriteLine("wrapped: p0,")
					b.WriteLine("timing: p1,")
					b.WriteLine("}")
				},
			},
		},
	}

	b.WriteFile(file)
}

func constructorName(iface types.Interface) string {
	return fmt.Sprintf("New%sTimer", strings.Title(iface.Name))
}

var contextNamedParam = types.NamedParam{Pkg: "context", Typ: "Context"}

func getFirstContextParamOffset(m types.Method) (int, bool) {
	for i, p := range m.Params {
		if np, ok := p.(types.NamedParam); ok {
			if np == contextNamedParam {
				return i, true
			}
		}
	}
	return -1, false
}
