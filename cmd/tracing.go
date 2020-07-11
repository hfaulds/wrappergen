package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hfaulds/wrappergen/gen/tracing"
)

type TracingFlags struct {
	InterfaceName string
}

func (c *TracingFlags) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.InterfaceName, "interface", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("interface")
}

func Tracing(rootFlags *RootFlags, tracingFlags *TracingFlags) error {
	rootConf, err := rootFlags.BuildConfig()
	if err != nil {
		return err
	}

	iface, ok := rootConf.Pkg.FindInterface(tracingFlags.InterfaceName)
	if !ok {
		return errors.New("Could not find interface")
	}

	if tracing.ShouldSkipInterface(iface) {
		return errors.New("Could not find any methods taking context")
	}

	tracing.Gen(rootConf.Builder, iface)

	rootConf.Builder.Flush("./gen_" + tracingFlags.InterfaceName + "Tracing.go")

	return nil
}
