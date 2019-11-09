package testdata

import i0 "context"
import i1 "bytes"
import trace "github.com/hfaulds/tracer/trace"

type traceanotherMethodsWithContext struct {
	wrapped anotherMethodsWithContext
}

func NewanotherMethodsWithContextTracer(p0 anotherMethodsWithContext) anotherMethodsWithContext {
	return traceanotherMethodsWithContext{
		wrapped: p0,
	}
}

func (t traceanotherMethodsWithContext) withContext(p0 i0.Context) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.withContext(ctx)
}

type tracemethodsWithContext struct {
	wrapped methodsWithContext
}

func NewmethodsWithContextTracer(p0 methodsWithContext) methodsWithContext {
	return tracemethodsWithContext{
		wrapped: p0,
	}
}

func (t tracemethodsWithContext) arrayType(p0 i0.Context, p1 [10]int) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.arrayType(ctx, p1)
}

func (t tracemethodsWithContext) interfaceType(p0 i0.Context, p1 interface {
	Foo(p0 string) int
},
) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.interfaceType(ctx, p1)
}

func (t tracemethodsWithContext) interfaceTypeWithEmbed(p0 i0.Context, p1 interface {
	Foo(p0 string) int
	withoutContext()
},
) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.interfaceTypeWithEmbed(ctx, p1)
}

func (t tracemethodsWithContext) mapType(p0 i0.Context, p1 map[int]string) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.mapType(ctx, p1)
}

func (t tracemethodsWithContext) namedAndBasicTypes(p0 i0.Context, p1 int, p2 i1.Buffer) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.namedAndBasicTypes(ctx, p1, p2)
}

func (t tracemethodsWithContext) pointerType(p0 i0.Context, p1 *int) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.pointerType(ctx, p1)
}

func (t tracemethodsWithContext) sliceType(p0 i0.Context, p1 []int) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.sliceType(ctx, p1)
}

func (t tracemethodsWithContext) withContext(p0 i0.Context) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.withContext(ctx)
}

func (t tracemethodsWithContext) withContextAsSecondArg(p0 int, p1 i0.Context) {
	ctx, span := trace.ChildSpan(p1)
	defer span.Close()
	t.wrapped.withContextAsSecondArg(p0, ctx)
}

func (t tracemethodsWithContext) withReturnType(p0 i0.Context) string {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	return t.wrapped.withReturnType(ctx)
}

func (t tracemethodsWithContext) withoutContext() {
	t.wrapped.withoutContext()
}
