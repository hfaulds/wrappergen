package tracing

import (
	"fmt"
	"io"
	"strings"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse/types"
)

func Gen(b io.Writer, iface types.Interface, importMap gen.ImportMap, tracePkg string) string {
	fmt.Fprintf(b, "import trace \"%s\"\n", tracePkg)

	tracingStruct := types.Struct{
		Name:  fmt.Sprintf("trace%s", iface.Name),
		Attrs: map[string]types.Param{"wrapped": types.NamedParam{Typ: iface.Name}},
	}

	gen.GenStruct(b, importMap, tracingStruct)

	tracingStructConstructor := types.Method{
		Name:    fmt.Sprintf("New%sTracer", strings.Title(iface.Name)),
		Params:  []types.Param{types.NamedParam{Typ: iface.Name}},
		Returns: []types.Param{types.NamedParam{Typ: iface.Name}},
	}

	gen.GenMethod(b, importMap, nil, tracingStructConstructor, func(b io.Writer) {
		fmt.Fprintf(b, "return %s{\n", tracingStruct.Name)
		fmt.Fprintf(b, "wrapped: p0,\n")
		fmt.Fprintf(b, "}\n")
	})

	for _, m := range iface.Methods {
		gen.GenMethod(b, importMap, &tracingStruct, m, func(b io.Writer) {
			// only add tracing if there a context
			offset, ok := getFirstContextParamOffset(m)
			if ok {
				fmt.Fprintf(b, "ctx, span := trace.ChildSpan(p%d, trace.OpName(\"%s\"))\n", offset, m.Name)
				fmt.Fprint(b, "defer span.Finish()\n")
			}
			generateWrappedCall(b, m, offset)
		})
	}
	fmt.Fprint(b, "\n")

	return tracingStructConstructor.Name
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

func generateWrappedCall(b io.Writer, m types.Method, contextOffset int) {
	fmt.Fprint(b, "")
	numReturns := len(m.Returns)
	errorOffset, returnsError := getLastErrorReturnOffset(m)
	if numReturns > 0 {
		if returnsError {
			for i := 0; i < numReturns; i++ {
				fmt.Fprintf(b, "r%d", i)
				if i != numReturns-1 {
					fmt.Fprint(b, ", ")
				}
			}
			fmt.Fprint(b, " := ")
		} else {
			fmt.Fprint(b, "return ")
		}
	}
	fmt.Fprintf(b, "t.wrapped.%s(", m.Name)
	for i := 0; i < len(m.Params); i++ {
		if i == contextOffset {
			fmt.Fprint(b, "ctx")
		} else {
			fmt.Fprintf(b, "p%d", i)
		}
		if i != len(m.Params)-1 {
			fmt.Fprint(b, ", ")
		}
	}
	fmt.Fprint(b, ")\n")
	if returnsError {
		fmt.Fprint(b, "return ")
		for i := 0; i < numReturns; i++ {
			if i == errorOffset {
				fmt.Fprintf(b, "span.WithError(r%d)", i)
			} else {
				fmt.Fprintf(b, "r%d", i)
			}
			if i != numReturns-1 {
				fmt.Fprint(b, ", ")
			}
		}
		fmt.Fprint(b, "\n")
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
