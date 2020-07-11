package tracing

import (
	"context"
)

type StartSpanOption func(interface{})

func OpName(opName string) func(interface{}) {
	return func(_ interface{}) {}
}

type Tracing interface {
	ChildSpan(context.Context, StartSpanOption)
}
