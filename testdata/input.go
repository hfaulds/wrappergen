package testdata

import (
	"bytes"
	"context"
)

type noMethods interface {
}

type noMethodsWithContext interface {
	withoutContext()
}

type methodsWithContext interface {
	withoutContext()
	withContext(context.Context)
	namedAndBasicTypes(context.Context, int, bytes.Buffer)
	arrayType(context.Context, [10]int)
	sliceType(context.Context, []int)
	pointerType(context.Context, *int)
	mapType(context.Context, map[int]string)
	withReturnType(context.Context) string
}
