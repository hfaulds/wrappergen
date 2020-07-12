package tracing

import (
	"context"
)

type Tracing interface {
	ChildSpan(context.Context, interface{}...)
	OpName(opName string) func(interface{})
}
