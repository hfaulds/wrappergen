package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse"
	"github.com/hfaulds/tracer/parse/types"
)

type RootFlags struct {
	Stdout bool
	Indir  string
}

func (f *RootFlags) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&f.Stdout, "stdout", false, "directory to write mocks to")
	cmd.PersistentFlags().StringVar(&f.Indir, "indir", ".", "directory to search for interface")
}

type RootConfig struct {
	Pkg     *types.Package
	Builder gen.Builder
}

func (f RootFlags) BuildConfig() (*RootConfig, error) {
	pkg, err := parse.ParseDir(f.Indir)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse %s")
	}

	conf := &RootConfig{
		Pkg:     pkg,
		Builder: gen.NewBuilder(pkg, f.Stdout),
	}
	return conf, nil
}
