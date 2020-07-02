package timing

import (
	"fmt"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

var timingType = types.NamedParam{
	Pkg: "github.com/hfaulds/tracer/timing",
	Typ: "Timing",
}

func Gen(b gen.Builder, iface types.Interface) string {
	timingStruct := types.Struct{
		Name: fmt.Sprintf("time%s", iface.Name),
		Attrs: []types.Var{
			{Name: "wrapped", Type: types.NamedParam{Typ: iface.Name}},
			{Name: "timing", Type: timingType},
		},
	}

	b.WriteStruct(timingStruct)

	timingStructConstructor := types.Method{
		Name: fmt.Sprintf("New%sTimer", strings.Title(iface.Name)),
		Params: []types.Param{
			types.NamedParam{Typ: iface.Name},
			timingType,
		},
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}

	b.WriteMethod(nil, timingStructConstructor, func(b gen.Builder) {
		b.WriteLine("return %s{", timingStruct.Name)
		b.WriteLine("wrapped: p0,")
		b.WriteLine("timing: p1,")
		b.WriteLine("}")
	})

	for _, m := range iface.Methods {
		b.WriteMethod(&timingStruct, m, func(b gen.Builder) {
			b.WriteLine("timer := t.timing.Timer()")
			b.WriteLine("defer timer.End(ctx, \"%s\")", m.Name)
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
		})
	}

	return timingStructConstructor.Name
}
