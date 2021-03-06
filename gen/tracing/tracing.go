package tracing

import (
	"fmt"
	"strings"

	"github.com/hfaulds/wrappergen/gen"
	"github.com/hfaulds/wrappergen/parse/types"
)

var traceType = types.NamedParam{
	Pkg: "github.com/hfaulds/wrappergen/tracing",
	Typ: "Tracing",
}

var TracingWrapper = func(iface types.Interface) gen.Wrapper {
	return gen.Wrapper{
		Constructor: constructorName(iface),
		Arguments:   []types.NamedParam{traceType},
	}
}

func Gen(b gen.Builder, iface types.Interface) {
	tracingStruct := gen.Struct{
		Struct: types.Struct{
			Name: fmt.Sprintf("trace%s", strings.Title(iface.Name)),
			Attrs: []types.Var{
				{Name: "wrapped", Type: types.NamedParam{Typ: iface.Name}},
				{Name: "trace", Type: traceType},
			},
		},
		Methods: make([]gen.Method, len(iface.Methods)),
	}

	for i, m := range iface.Methods {
		tracingStruct.Methods[i] = gen.Method{
			Method: m,
			Callback: func(b gen.Builder, m types.Method) {
				// only add tracing if there a context
				offset, ok := getFirstContextParamOffset(m)
				if ok {
					b.WriteLine("ctx, span := t.trace.ChildSpan(p%d, \"%s\")", offset, m.Name)
					b.WriteLine("defer span.Finish()")
				}
				generateWrappedCall(b, m, offset)
			},
		}
	}

	file := gen.File{
		Structs: []gen.Struct{tracingStruct},
		Methods: []gen.Method{
			{
				Method: types.Method{
					Name: constructorName(iface),
					Params: []types.Param{
						types.NamedParam{Typ: iface.Name},
						traceType,
					},
					Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
				},
				Callback: func(b gen.Builder, _ types.Method) {
					b.WriteLine("return %s{", tracingStruct.Name)
					b.WriteLine("wrapped: p0,")
					b.WriteLine("trace: p1,")
					b.WriteLine("}")
				},
			},
		},
	}
	b.WriteFile(file)
}

func ShouldSkipInterface(i types.Interface) bool {
	for _, m := range i.Methods {
		if _, ok := getFirstContextParamOffset(m); ok {
			return false
		}
	}
	return true
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

func generateWrappedCall(b gen.Builder, m types.Method, contextOffset int) {
	numReturns := len(m.Returns)
	errorOffset, returnsError := getLastErrorReturnOffset(m)
	if numReturns > 0 {
		if returnsError {
			for i := 0; i < numReturns; i++ {
				b.Write("r%d", i)
				if i != numReturns-1 {
					b.Write(", ")
				}
			}
			b.Write(" := ")
		} else {
			b.Write("return ")
		}
	}
	b.Write("t.wrapped.%s(", m.Name)
	for i := 0; i < len(m.Params); i++ {
		if i == contextOffset {
			b.Write("ctx")
		} else {
			b.Write("p%d", i)
		}
		if i != len(m.Params)-1 {
			b.Write(", ")
		}
	}
	b.WriteLine(")")
	if returnsError {
		b.Write("return ")
		for i := 0; i < numReturns; i++ {
			if i == errorOffset {
				b.Write("span.WithError(r%d)", i)
			} else {
				b.Write("r%d", i)
			}
			if i != numReturns-1 {
				b.Write(", ")
			}
		}
		b.WriteLine("")
	}
}

var errorNamedParam = types.NamedParam{Pkg: "", Typ: "error"}

func getLastErrorReturnOffset(m types.Method) (int, bool) {
	for i := len(m.Returns) - 1; i >= 0; i-- {
		p := m.Returns[i]
		if np, ok := p.(types.NamedParam); ok {
			if np == errorNamedParam {
				return i, true
			}
		}
	}
	return -1, false
}

func constructorName(iface types.Interface) string {
	return fmt.Sprintf("New%sTracer", strings.Title(iface.Name))
}
