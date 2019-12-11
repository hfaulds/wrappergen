package timing

import (
	"fmt"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b *gen.Builder, iface types.Interface, importMap gen.ImportMap, timingAttr string) string {
	timingStruct := types.Struct{
		Name:  fmt.Sprintf("time%s", iface.Name),
		Attrs: map[string]types.Param{"wrapped": types.NamedParam{Typ: iface.Name}},
	}

	b.WriteStruct(importMap, timingStruct)

	timingStructConstructor := types.Method{
		Name:    fmt.Sprintf("New%sTimer", strings.Title(iface.Name)),
		Params:  []types.Param{types.NamedParam{Typ: iface.Name}},
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}

	b.WriteMethod(importMap, nil, timingStructConstructor, func(b *gen.Builder) {
		fmt.Fprintf(b, "return %s{\n", timingStruct.Name)
		fmt.Fprintf(b, "wrapped: p0,\n")
		fmt.Fprintf(b, "}\n")
	})

	for _, m := range iface.Methods {
		b.WriteMethod(importMap, &timingStruct, m, func(b *gen.Builder) {
			fmt.Fprintf(b, "timer := t.%s.Timer()\n", timingAttr)
			fmt.Fprintf(b, "defer timer.End(ctx, \"%s\")\n", m.Name)
			numReturns := len(m.Returns)
			if numReturns > 0 {
				fmt.Fprint(b, "return ")
			}
			fmt.Fprintf(b, "t.wrapped.%s(", m.Name)
			for i := 0; i < len(m.Params); i++ {
				fmt.Fprintf(b, "p%d", i)
				if i != len(m.Params)-1 {
					fmt.Fprint(b, ", ")
				}
			}
			fmt.Fprint(b, ")")
		})
	}

	return timingStructConstructor.Name
}

func StructHasTimingAttr(strct types.Struct, timingAttr string) bool {
	for attrName := range strct.Attrs {
		if attrName == timingAttr {
			return true
		}
	}
	return false
}