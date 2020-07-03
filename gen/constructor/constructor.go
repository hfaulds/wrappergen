package constructor

import (
	"fmt"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b gen.Builder, iface types.Interface, strct types.Struct, wrappers []gen.Wrapper) {
	for w, wrapper := range wrappers {
		for a, arg := range wrapper.Arguments {
			b.AddImport(fmt.Sprintf("w%da%d", w, a), arg.Pkg)
		}
	}
	b.WriteImports()

	attrTypes := make([]types.Param, 0, len(strct.Attrs))
	for _, attr := range strct.Attrs {
		attrTypes = append(attrTypes, attr.Type)
	}

	for _, wrapper := range wrappers {
		for _, arg := range wrapper.Arguments {
			attrTypes = append(attrTypes, arg)
		}
	}

	b.WriteMethod(nil, types.Method{
		Name:    fmt.Sprintf("New%s", strings.Title(iface.Name)),
		Params:  attrTypes,
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}, func(b gen.Builder) {
		b.Write("return ")
		for _, wrapper := range wrappers {
			b.Write("%s(\n", wrapper.Constructor)
		}

		b.WriteLine("%s{", strct.Name)
		argIndex := 0
		for _, attr := range strct.Attrs {
			b.Write("%s: p%d,\n", attr.Name, argIndex)
			argIndex++
		}
		b.Write("}")

		if len(attrTypes) > len(strct.Attrs) {
			b.Write(",")
		}
		b.Write("\n")

		for i, wrapper := range wrappers {
			for range wrapper.Arguments {
				b.Write("p%d,\n", argIndex)
				argIndex++
			}
			b.Write(")")
			if i < len(wrappers) - 1 {
				b.Write(",\n")
			}
		}

	})
}
