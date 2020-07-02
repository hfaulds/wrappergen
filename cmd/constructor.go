package cmd

import (
	"github.com/spf13/cobra"
)

type ConstructorConfig struct {
	StructName    string
	InterfaceName string
}

func (c *ConstructorConfig) Init(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.InterfaceName, "interface", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("interface")
	cmd.Flags().StringVar(&c.InterfaceName, "struct", "", "name of interface to generate wrappers for")
	cmd.MarkFlagRequired("struct")
}
