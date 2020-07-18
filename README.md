# Wrappergen

**Warning** Very WIP, no tests, no guarantees.

Code-gen to add tracing, timing wrappers and constructors.

Timing and tracing code is very common boilerplate scattered around entire codebases. Wrappergen aims make use of codegen to remove that boilerplate.

### Example

This repo contains `./example/example.go` which looks like so:

```
package example

import (
	"context"
)

type Example interface {
	Test(context.Context, int64) uint8
}

type example struct {
	attribute string
}

func (e *example) Test(ctx context.Context, i int64) uint8 {
	return 5
}
```

#### Timing

```
$ wrappergen timing --indir ./example/ --interface Example
```

```
// Code generated by wrappergen v0.0.1. DO NOT EDIT.
package example

import i0 "github.com/hfaulds/wrappergen/timing"
import i1 "context"

func NewExampleTimer(p0 Example, p1 i0.Timing) Example {
	return timeExample{
		wrapped: p0,
		timing:  p1,
	}
}

type timeExample struct {
	wrapped Example
	timing  i0.Timing
}

func (t timeExample) Test(p0 i1.Context, p1 int64) uint8 {
	timer := t.timing.Timer()
	defer timer.End(ctx, "Test")
	return t.wrapped.Test(p0, p1)
}
```

#### Tracing

```
$ wrappergen tracing --indir ./example/ --interface Example
```

```
// Code generated by wrappergen v0.0.1. DO NOT EDIT.
package example

import i0 "github.com/hfaulds/wrappergen/tracing"
import i1 "context"

func NewExampleTracer(p0 Example, p1 i0.Tracing) Example {
	return traceExample{
		wrapped: p0,
		trace:   p1,
	}
}

type traceExample struct {
	wrapped Example
	trace   i0.Tracing
}

func (t traceExample) Test(p0 i1.Context, p1 int64) uint8 {
	ctx, span := t.trace.ChildSpan(p0, "Test")
	defer span.Finish()
	return t.wrapped.Test(ctx, p1)
}
```


#### Constructor

```
$ wrappergen constructor --indir ./example/ --interface Example --struct example --timing --tracing
```

```
// Code generated by wrappergen v0.0.1. DO NOT EDIT.
package example

import i0 "github.com/hfaulds/wrappergen/tracing"
import i1 "github.com/hfaulds/wrappergen/timing"

func NewExample(p0 string, p1 i0.Tracing, p2 i1.Timing) Example {
	return NewExampleTracer(
		NewExampleTimer(
			example{
				attribute: p0,
			},
			p1,
		),
		p2,
	)
}
```

### TODO

- [ ] Testing
- [ ] CI
- [ ] Tidy up code generation
