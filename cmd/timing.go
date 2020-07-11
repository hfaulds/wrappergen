package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hfaulds/wrappergen/gen/timing"
)

type TimingFlags struct {
	InterfaceName string
}

func (f *TimingFlags) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&f.InterfaceName, "interface", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("interface")
}

func Timing(rootFlags *RootFlags, timingFlags *TimingFlags) error {
	rootConf, err := rootFlags.BuildConfig()
	if err != nil {
		return err
	}

	iface, ok := rootConf.Pkg.FindInterface(timingFlags.InterfaceName)
	if !ok {
		return errors.New("Could not find interface")
	}

	timing.Gen(rootConf.Builder, iface)

	rootConf.Builder.Flush("./gen_" + timingFlags.InterfaceName + "Timing.go")

	return nil
}
