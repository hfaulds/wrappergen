package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/parse"
	"github.com/hfaulds/tracer/parse/types"
)

type RootFlags struct {
	Stdout bool
	Outdir string
	Outpkg string
	Inpkg  bool
	Indir  string
}

func (f *RootFlags) Init(cmd *cobra.Command) {
	cmd.PersistentFlags().BoolVar(&f.Stdout, "stdout", false, "directory to write mocks to")
	cmd.PersistentFlags().StringVar(&f.Outdir, "outdir", ".", "directory to write mocks to")
	cmd.PersistentFlags().StringVar(&f.Outpkg, "outpkg", "tracing", "name of generated package")
	cmd.PersistentFlags().BoolVar(&f.Inpkg, "inpkg", false, "generate a mock that goes inside the original package")
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

	dst := os.Stdout
	if !f.Stdout {
		if err := os.MkdirAll(filepath.Dir(f.Outdir), os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "Unable to create directory: %v")
		}
		f, err := os.Create(f.Outdir)
		if err != nil {
			return nil, errors.Wrap(err, "Failed opening destination file: %v")
		}
		defer f.Close()
		dst = f
	}

	builder := gen.NewBuilder(pkg, dst)

	conf := &RootConfig{
		Pkg:     pkg,
		Builder: builder,
	}
	return conf, nil
}
