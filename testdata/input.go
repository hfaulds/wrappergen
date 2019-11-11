package testdata

import (
	"bytes"
	"context"
	"io"
)

type internalType struct {
	foo string
}

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
	internalTypeParam(context.Context, internalType)
	arrayType(context.Context, [10]int)
	sliceType(context.Context, []int)
	pointerType(context.Context, *int)
	mapType(context.Context, map[int]string)
	returnBasicType(context.Context) string
	returnNamedType(context.Context) io.Reader
	returnInternalType(context.Context) internalType
	interfaceType(context.Context, interface{ Foo(string) int })
	interfaceTypeWithEmbed(context.Context, interface {
		noMethodsWithContext
		Foo(string) int
	})
}

type anotherMethodsWithContext interface {
	withContext(context.Context)
}
