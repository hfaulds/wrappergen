package timing

import (
	"fmt"
	"io"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b io.Writer, iface types.Interface, importMap gen.ImportMap, timingAttr string) string {
	timingStruct := types.Struct{
		Name:  fmt.Sprintf("time%s", iface.Name),
		Attrs: map[string]types.Param{"wrapped": types.NamedParam{Typ: iface.Name}},
	}

	gen.GenStruct(b, importMap, timingStruct)

	timingStructConstructor := types.Method{
		Name:    fmt.Sprintf("New%sTimer", strings.Title(iface.Name)),
		Params:  []types.Param{types.NamedParam{Typ: iface.Name}},
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}

	gen.GenMethod(b, importMap, nil, timingStructConstructor, func(b io.Writer) {
		fmt.Fprintf(b, "\treturn %s{\n", timingStruct.Name)
		fmt.Fprintf(b, "\t\twrapped: p0,\n")
		fmt.Fprintf(b, "\t}\n")
	})

	for _, m := range iface.Methods {
		gen.GenMethod(b, importMap, &timingStruct, m, func(b io.Writer) {
			fmt.Fprintf(b, "\ttimer := t.%s.Timer()\n", timingAttr)
			fmt.Fprintf(b, "\tdefer timer.End(ctx, \"%s\")\n\t", m.Name)
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
			fmt.Fprint(b, ")\n")
		})
	}
	fmt.Fprint(b, "\n")

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
