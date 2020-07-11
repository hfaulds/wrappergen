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
