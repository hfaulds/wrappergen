package example

import (
	"context"
)

//go:generate go run ../main.go timing --interface Example
//go:generate go run ../main.go tracing --interface Example
//go:generate go run ../main.go constructor --interface Example --struct example --timing --tracing
type Example interface {
	Test(context.Context, int64) error
}

type example struct {
	attribute string
}

func (e *example) Test(ctx context.Context, i int64) error {
	return nil
}
