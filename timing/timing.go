package timing

import (
	"context"
)

type Timer interface {
	End(context.Context, string)
}

type Timing interface {
	Timer() Timer
}
