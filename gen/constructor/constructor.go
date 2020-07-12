package constructor

import (
	"fmt"
	"strings"

	"github.com/hfaulds/wrappergen/gen"
	"github.com/hfaulds/wrappergen/parse/types"
)

func Gen(b gen.Builder, iface types.Interface, strct types.Struct, wrappers []gen.Wrapper) {
	attrTypes := make([]types.Param, 0, len(strct.Attrs))
	for _, attr := range strct.Attrs {
		attrTypes = append(attrTypes, attr.Type)
	}

	for _, wrapper := range wrappers {
		for _, arg := range wrapper.Arguments {
			attrTypes = append(attrTypes, arg)
		}
	}

	file := gen.File{
		Structs: []gen.Struct{},
		Methods: []gen.Method{
			{
				Method: types.Method{
					Name:    fmt.Sprintf("New%s", strings.Title(iface.Name)),
					Params:  attrTypes,
					Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
				},
				Callback: func(b gen.Builder, _ types.Method) {
					b.Write("return ")
					for i := len(wrappers) - 1; i >= 0; i-- {
						wrapper := wrappers[i]
						b.Write("%s(\n", wrapper.Constructor)
					}

					b.WriteLine("&%s{", strct.Name)
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
						if i < len(wrappers)-1 {
							b.Write(",\n")
						}
					}
				},
			},
		},
	}

	b.WriteFile(file)
}
