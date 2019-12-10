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
	returnMultipleErrors(context.Context) (error, error)
	interfaceType(context.Context, interface{ Foo(string) int })
	interfaceTypeEmty(context.Context, interface{})
	interfaceTypeWithEmbed(context.Context, interface {
		noMethodsWithContext
		Foo(string) int
	})
}

type anotherMethodsWithContext interface {
	withContext(context.Context)
}

type A struct {
	B string
}

func (a A) withoutContext()                                              {}
func (a A) withContext(context.Context)                                  {}
func (a A) withContextAsSecondArg(int, context.Context)                  {}
func (a A) namedAndBasicTypes(context.Context, int, bytes.Buffer, error) {}
func (a A) internalTypeParam(context.Context, internalType)              {}
func (a A) arrayType(context.Context, [10]int)                           {}
func (a A) sliceType(context.Context, []int)                             {}
func (a A) pointerType(context.Context, *int)                            {}
func (a A) mapType(context.Context, map[int]string)                      {}
func (a A) returnNamedAndBasicTypes(context.Context) (string, io.Reader, error) {
	return "", nil, nil
}
func (a A) returnInternalType(context.Context) internalType {
	return internalType{}
}
func (a A) returnMultipleErrors(context.Context) (error, error) {
	return nil, nil
}
func (a A) interfaceType(context.Context, interface{ Foo(string) int }) {}
func (a A) interfaceTypeEmty(context.Context, interface{})              {}
func (a A) interfaceTypeWithEmbed(context.Context, interface {
	noMethodsWithContext
	Foo(string) int
}) {
}
