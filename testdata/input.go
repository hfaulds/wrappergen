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
	namedAndBasicTypes(context.Context, int, bytes.Buffer, error)
	internalTypeParam(context.Context, internalType)
	arrayType(context.Context, [10]int)
	sliceType(context.Context, []int)
	pointerType(context.Context, *int)
	mapType(context.Context, map[int]string)
	returnNamedAndBasicTypes(context.Context) (string, io.Reader, error)
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
