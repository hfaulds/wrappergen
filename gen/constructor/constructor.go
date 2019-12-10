package constructor

import (
	"fmt"
	"io"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b io.Writer, importMap gen.ImportMap, iface types.Interface, strct types.Struct, wrappers []string) {
	attrTypes := make([]types.Param, 0, len(strct.Attrs))
	for _, typ := range strct.Attrs {
		attrTypes = append(attrTypes, typ)
	}

	gen.GenMethod(b, importMap, nil, types.Method{
		Name:    fmt.Sprintf("New%s", strings.Title(iface.Name)),
		Params:  attrTypes,
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}, func(b io.Writer) {
		fmt.Fprintf(b, "\t return ")
		for _, wrapper := range wrappers {
			fmt.Fprintf(b, "%s(", wrapper)
		}

		fmt.Fprintf(b, "%s{\n", strct.Name)
		i := 0
		for name, _ := range strct.Attrs {
			fmt.Fprintf(b, "\t\t %s: a%d,\n", name, i)
			i++
		}
		fmt.Fprint(b, "\t }")

		for i := 0; i < len(wrappers); i++ {
			fmt.Fprint(b, ")")
		}
	})
}
