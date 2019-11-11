package trace

import "context"

type span interface {
	Finish()
	WithError(error) error
}

type opName string

func ChildSpan(ctx context.Context, _ opName) (context.Context, span) {
	return ctx, nil
}

func OpName(str string) opName {
	return opName(str)
}
