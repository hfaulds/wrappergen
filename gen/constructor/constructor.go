package constructor

import (
	"fmt"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b gen.Builder, iface types.Interface, strct types.Struct, wrappers []string) {
	attrTypes := make([]types.Param, 0, len(strct.Attrs))
	for _, attr := range strct.Attrs {
		attrTypes = append(attrTypes, attr.Type)
	}

	b.WriteMethod(nil, types.Method{
		Name:    fmt.Sprintf("New%s", strings.Title(iface.Name)),
		Params:  attrTypes,
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}, func(b gen.Builder) {
		b.Write("return ")
		for _, wrapper := range wrappers {
			b.Write("%s(", wrapper)
		}

		b.WriteLine("%s{", strct.Name)
		i := 0
		for _, attr := range strct.Attrs {
			b.Write("%s: a%d,\n", attr.Name, i)
			i++
		}
		b.Write("}")

		for i := 0; i < len(wrappers); i++ {
			b.Write(")")
		}
	})
}
