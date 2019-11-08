package test

import i0 "context"
import i1 "bytes"

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
	t.wrapped.arrayType(p0, p1)
}

func (t tracemethodsWithContext) mapType(p0 i0.Context, p1 map[int]string) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.mapType(p0, p1)
}

func (t tracemethodsWithContext) namedAndBasicTypes(p0 i0.Context, p1 int, p2 i1.Buffer) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.namedAndBasicTypes(p0, p1, p2)
}

func (t tracemethodsWithContext) pointerType(p0 i0.Context, p1 *int) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.pointerType(p0, p1)
}

func (t tracemethodsWithContext) sliceType(p0 i0.Context, p1 []int) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.sliceType(p0, p1)
}

func (t tracemethodsWithContext) withContext(p0 i0.Context) {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	t.wrapped.withContext(p0)
}

func (t tracemethodsWithContext) withReturnType(p0 i0.Context) string {
	ctx, span := trace.ChildSpan(p0)
	defer span.Close()
	return t.wrapped.withReturnType(p0)
}

func (t tracemethodsWithContext) withoutContext() {
	t.wrapped.withoutContext()
}
