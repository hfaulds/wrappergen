package constructor

import (
	"fmt"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b *gen.Builder, importMap gen.ImportMap, iface types.Interface, strct types.Struct, wrappers []string) {
	attrTypes := make([]types.Param, 0, len(strct.Attrs))
	for _, typ := range strct.Attrs {
		attrTypes = append(attrTypes, typ)
	}

	b.WriteMethod(importMap, nil, types.Method{
		Name:    fmt.Sprintf("New%s", strings.Title(iface.Name)),
		Params:  attrTypes,
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}, func(b *gen.Builder) {
		fmt.Fprintf(b, "return ")
		for _, wrapper := range wrappers {
			fmt.Fprintf(b, "%s(", wrapper)
		}

		fmt.Fprintf(b, "%s{\n", strct.Name)
		i := 0
		for name, _ := range strct.Attrs {
			fmt.Fprintf(b, "%s: a%d,\n", name, i)
			i++
		}
		fmt.Fprint(b, "}")

		for i := 0; i < len(wrappers); i++ {
			fmt.Fprint(b, ")")
		}
	})
}
