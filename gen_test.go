package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hfaulds/tracer/testdata"
)

/* TODO
test package resolving for each type
*/

func TestParseAndGen(t *testing.T) {
	pkg, err := ParseDir("./testdata")
	require.NoError(t, err)
	t.Log(pkg)
	tracePkg := "github.com/hfaulds/tracer/testdata/trace"
	generated := Generate(pkg, tracePkg)
	expected, err := ioutil.ReadFile("./testdata/trace.go")
	require.NoError(t, err)
	assert.Equal(t, string(expected), generated)
	assert.NotNil(t, testdata.NewMethodsWithContextTracer(nil))
}
