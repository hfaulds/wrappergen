package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAndGen(t *testing.T) {
	pkg, err := ParseDir("./testinput")
	require.NoError(t, err)
	t.Log(pkg)
	generated := Generate(pkg)
	expected, err := ioutil.ReadFile("./testoutput/output.go")
	require.NoError(t, err)
	assert.Equal(t, string(expected), generated)
}
