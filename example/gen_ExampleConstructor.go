// Code generated by wrappergen v0.0.1. DO NOT EDIT.
package example

import i0 "github.com/hfaulds/wrappergen/tracing"
import i1 "github.com/hfaulds/wrappergen/timing"

func NewExample(p0 string, p1 i0.Tracing, p2 i1.Timing) Example {
	return NewExampleTimer(
		NewExampleTracer(
			&example{
				attribute: p0,
			},
			p1,
		),
		p2,
	)
}
