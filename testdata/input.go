package testdata

import (
	"bytes"
	"context"
)

type returnType string

type noMethods interface {
}

type noMethodsWithContext interface {
	withoutContext()
}

type methodsWithContext interface {
	withoutContext()
	withContext(context.Context)
	withContextAsSecondArg(int, context.Context)
	namedAndBasicTypes(context.Context, int, bytes.Buffer)
	arrayType(context.Context, [10]int)
	sliceType(context.Context, []int)
	pointerType(context.Context, *int)
	mapType(context.Context, map[int]string)
	withReturnType(context.Context) string
	withInternalReturnType(context.Context) returnType
	interfaceType(context.Context, interface{ Foo(string) int })
	interfaceTypeWithEmbed(context.Context, interface {
		noMethodsWithContext
		Foo(string) int
	})
}

type anotherMethodsWithContext interface {
	withContext(context.Context)
}
