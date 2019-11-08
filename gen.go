package main

import (
	"fmt"
	"strings"
)

func generate(interfaces []Interface) string {
	var b strings.Builder
	for _, iface := range interfaces {
		structName := fmt.Sprintf("trace%s", iface.Name)
		fmt.Fprintf(&b, "type %s struct {\n", structName)
		fmt.Fprintf(&b, "	wrapped %s\n", iface.Name)
		fmt.Fprintf(&b, "}\n\n")

		generateMethodSig(&b, "", fmt.Sprintf("New%sTracer", iface.Name), []string{iface.Name}, []string{iface.Name})
		fmt.Fprintf(&b, "	return %s {\n", structName)
		fmt.Fprintf(&b, "		wrapped: %s,\n", iface.Name)
		fmt.Fprintf(&b, "	}\n")
		fmt.Fprintf(&b, "}\n\n")

		for _, method := range iface.Methods {
			generateMethodSig(&b, structName, method.Name, []string{}, []string{})
			fmt.Fprint(&b, "	ctx, span := trace.ChildSpan(ctx)\n")
			fmt.Fprint(&b, "	defer span.Close()\n")
			fmt.Fprintf(&b, "	return	wrapped.%s()\n", method.Name)
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
