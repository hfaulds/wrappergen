package trace

import "context"

func ChildSpan(ctx context.Context) (context.Context, span) {
	return ctx, span{}
}

type span struct {
}

func (s span) Close() {
}
