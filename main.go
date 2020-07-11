package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/hfaulds/wrappergen/cmd"
)

//go:generate gen tracing -interface=Client
//go:generate gen timing -interface=Client
//go:generate gen constructor -interface=Client -struct=client

var (
	rootFlags = &cmd.RootFlags{}
	rootCmd   = &cobra.Command{
		Use:   "gen",
		Short: "",
		Args:  cobra.NoArgs,
	}

	tracingFlags = &cmd.TracingFlags{}
	tracingCmd   = &cobra.Command{
		Use:   "tracing",
		Short: "",
		Args:  cobra.NoArgs,
		Run: runCommand(func() error {
			return cmd.Tracing(rootFlags, tracingFlags)
		}),
	}

	timingFlags = &cmd.TimingFlags{}
	timingCmd   = &cobra.Command{
		Use:   "timing",
		Short: "",
		Args:  cobra.NoArgs,
		Run: runCommand(func() error {
			return cmd.Timing(rootFlags, timingFlags)
		}),
	}

	constructorFlags = &cmd.ConstructorFlags{}
	constructorCmd   = &cobra.Command{
		Use:   "constructor",
		Short: "",
		Args:  cobra.NoArgs,
		Run: runCommand(func() error {
			return cmd.Constructor(rootFlags, constructorFlags)
		}),
	}
)

func runCommand(fn func() error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		err := fn()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func main() {
	rootFlags.Init(rootCmd)

	rootCmd.AddCommand(tracingCmd)
	tracingFlags.Init(tracingCmd)

	rootCmd.AddCommand(timingCmd)
	timingFlags.Init(timingCmd)

	rootCmd.AddCommand(constructorCmd)
	constructorFlags.Init(constructorCmd)

	rootCmd.Execute()
}
