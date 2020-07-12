package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/hfaulds/wrappergen/gen"
	"github.com/hfaulds/wrappergen/gen/constructor"
	"github.com/hfaulds/wrappergen/gen/timing"
	"github.com/hfaulds/wrappergen/gen/tracing"
)

type ConstructorFlags struct {
	StructName    string
	InterfaceName string
	Tracing       bool
	Timing        bool
}

func (c *ConstructorFlags) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.InterfaceName, "interface", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("interface")
	cmd.Flags().StringVar(&c.StructName, "struct", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("struct")
	cmd.Flags().BoolVar(&c.Tracing, "tracing", false, "whether to create a timing wrapper")
	cmd.Flags().BoolVar(&c.Timing, "timing", false, "whether to create a tracing wrapper")
}

func Constructor(rootFlags *RootFlags, constructorFlags *ConstructorFlags) error {
	rootConf, err := rootFlags.BuildConfig()
	if err != nil {
		return err
	}

	iface, ok := rootConf.Pkg.FindInterface(constructorFlags.InterfaceName)
	if !ok {
		return errors.New("Could not find interface")
	}

	strct, ok := rootConf.Pkg.FindStruct(constructorFlags.StructName)
	if !ok {
		return errors.New("Could not find struct")
	}

	wrappers := []gen.Wrapper{}
	if constructorFlags.Tracing {
		wrappers = append(wrappers, tracing.TracingWrapper(iface))
	}
	if constructorFlags.Timing {
		wrappers = append(wrappers, timing.TimingWrapper(iface))
	}

	constructor.Gen(rootConf.Builder, iface, strct, wrappers)

	return rootConf.Builder.Flush("./gen_" + constructorFlags.InterfaceName + "Constructor.go")
}
