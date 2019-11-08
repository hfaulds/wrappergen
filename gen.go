package main

import (
	"fmt"
	"strings"
)

/*
TODO:
imports
method arguments
*/

func generate(interfaces []Interface) string {
	var b strings.Builder
	for _, iface := range interfaces {
		// Skip interfaces where no methods have context.Context as an argument
		if shouldSkipInterface(iface) {
			continue
		}

		/*
			type traceExample struct {
				wrapped Example
			}
		*/
		structName := fmt.Sprintf("trace%s", iface.Name)
		fmt.Fprintf(&b, "type %s struct {\n", structName)
		fmt.Fprintf(&b, "	wrapped %s\n", iface.Name)
		fmt.Fprintf(&b, "}\n\n")

		/* func NewExampleTracer(p0 Example) Example {
			return traceExample {
				wrapped p0,
			}
		}*/
		generateMethodSig(&b, "", fmt.Sprintf("New%sTracer", iface.Name), []string{iface.Name}, []string{iface.Name})
		fmt.Fprintf(&b, "\treturn %s {\n", structName)
		fmt.Fprintf(&b, "\t\twrapped: p0,\n")
		fmt.Fprintf(&b, "\t}\n")
		fmt.Fprintf(&b, "}\n\n")

		/*
			func (t traceExample) Foo(p0 context.Context, p1) error {
				ctx, span := trace.ChildSpan(p0)
				defer span.Close()
				return t.wrapped(p0,p1)
			}
		*/
		for i, m := range iface.Methods {
			generateMethodSig(&b, structName, m.Name, []string{}, []string{})
			// only add tracing if there a context
			offset, ok := getFirstContextParamOffset(m)
			if ok && i == offset {
				fmt.Fprintf(&b, "\tctx, span := trace.ChildSpan(p%d)\n", offset)
				fmt.Fprint(&b, "\tdefer span.Close()\n")
			}
			fmt.Fprintf(&b, "\treturn t.wrapped.%s()\n", m.Name)
			fmt.Fprintf(&b, "}\n\n")
		}
	}
	return b.String()
}

func generateMethodSig(b *strings.Builder, implementor, methodName string, params, returns []string) {
	fmt.Fprint(b, "func ")
	if implementor != "" {
		fmt.Fprintf(b, "(t %s) ", implementor)
	}
	fmt.Fprintf(b, "%s(", methodName)
	for i, param := range params {
		fmt.Fprintf(b, "p%d %s", i, param)
		if i < len(params)-1 {
			fmt.Fprint(b, ", ")
		}
	}
	fmt.Fprint(b, ") ")
	if len(returns) > 1 {
		fmt.Fprint(b, "(")
	}
	for i, r := range returns {
		fmt.Fprint(b, r)
		if i < len(params)-1 {
			fmt.Fprint(b, ", ")
		}
	}
	if len(returns) > 1 {
		fmt.Fprint(b, ")")
	}
	fmt.Fprint(b, " {\n")
}

func shouldSkipInterface(i Interface) bool {
	for _, m := range i.Methods {
		if _, ok := getFirstContextParamOffset(m); ok {
			return false
		}
	}
	return true
}

func getFirstContextParamOffset(m method) (int, bool) {
	for i, p := range m.Params {
		if p.Type == "context.Context" {
			return i, true
		}
	}
	return 0, false
}
