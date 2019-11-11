package testdata

import i0 "context"
import i1 "bytes"
import i2 "io"

type traceanotherMethodsWithContext struct {
	wrapped   anotherMethodsWithContext
	childSpan func(i0.Context) (i0.Context, interface{ Close() })
}

func NewanotherMethodsWithContextTracer(p0 anotherMethodsWithContext, p1 func(i0.Context) (i0.Context, interface{ Close() })) anotherMethodsWithContext {
	return traceanotherMethodsWithContext{
		wrapped:   p0,
		childSpan: p1,
	}
}

func (t traceanotherMethodsWithContext) withContext(p0 i0.Context) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.withContext(ctx)
}

type tracemethodsWithContext struct {
	wrapped   methodsWithContext
	childSpan func(i0.Context) (i0.Context, interface{ Close() })
}

func NewmethodsWithContextTracer(p0 methodsWithContext, p1 func(i0.Context) (i0.Context, interface{ Close() })) methodsWithContext {
	return tracemethodsWithContext{
		wrapped:   p0,
		childSpan: p1,
	}
}

func (t tracemethodsWithContext) arrayType(p0 i0.Context, p1 [10]int) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.arrayType(ctx, p1)
}

func (t tracemethodsWithContext) interfaceType(p0 i0.Context, p1 interface {
	Foo(p0 string) int
},
) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.interfaceType(ctx, p1)
}

func (t tracemethodsWithContext) interfaceTypeWithEmbed(p0 i0.Context, p1 interface {
	Foo(p0 string) int
	withoutContext()
},
) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.interfaceTypeWithEmbed(ctx, p1)
}

func (t tracemethodsWithContext) internalTypeParam(p0 i0.Context, p1 internalType) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.internalTypeParam(ctx, p1)
}

func (t tracemethodsWithContext) mapType(p0 i0.Context, p1 map[int]string) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.mapType(ctx, p1)
}

func (t tracemethodsWithContext) namedAndBasicTypes(p0 i0.Context, p1 int, p2 i1.Buffer) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.namedAndBasicTypes(ctx, p1, p2)
}

func (t tracemethodsWithContext) pointerType(p0 i0.Context, p1 *int) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.pointerType(ctx, p1)
}

func (t tracemethodsWithContext) returnBasicType(p0 i0.Context) string {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	return t.wrapped.returnBasicType(ctx)
}

func (t tracemethodsWithContext) returnInternalType(p0 i0.Context) internalType {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	return t.wrapped.returnInternalType(ctx)
}

func (t tracemethodsWithContext) returnNamedType(p0 i0.Context) i2.Reader {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	return t.wrapped.returnNamedType(ctx)
}

func (t tracemethodsWithContext) sliceType(p0 i0.Context, p1 []int) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.sliceType(ctx, p1)
}

func (t tracemethodsWithContext) withContext(p0 i0.Context) {
	ctx, span := t.childSpan(p0)
	defer span.Close()
	t.wrapped.withContext(ctx)
}

func (t tracemethodsWithContext) withContextAsSecondArg(p0 int, p1 i0.Context) {
	ctx, span := t.childSpan(p1)
	defer span.Close()
	t.wrapped.withContextAsSecondArg(p0, ctx)
}

func (t tracemethodsWithContext) withoutContext() {
	t.wrapped.withoutContext()
}
