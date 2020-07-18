package tracing

import (
	"context"
)

type Tracing interface {
	ChildSpan(context.Context, string) (context.Context, Span)
}

type Span interface {
	WithError(error) error
	Finish()
}
