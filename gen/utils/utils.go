package utils

import (
	"fmt"
	"io"

	"github.com/hfaulds/tracer/parse/types"
)

//func WriteStruct(b io.Writer, strct types.Struct) {
//fmt.Fprintf(&b, "type %s struct {\n", strct.Name)
//for name, typ := range strct.Attrs {
//fmt.Fprintf(&b, "\t%s", name)
//WriteParam(b, typ)
//fmt.Fprint("%s\n")
//}
//fmt.Fprintf(&b, "}")
//}
