package tracing

import (
	"context"
)

type Tracing interface {
	ChildSpan(context.Context, ...interface{}) (context.Context, Span)
	OpName(opName string) func(interface{})
}

type Span interface {
	Finish()
}
