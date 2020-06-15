package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/hfaulds/tracer/gen"
	"github.com/hfaulds/tracer/gen/constructor"
	"github.com/hfaulds/tracer/gen/timing"
	"github.com/hfaulds/tracer/gen/tracing"
	"github.com/hfaulds/tracer/flags"
	"github.com/hfaulds/tracer/parse"
	"github.com/hfaulds/tracer/parse/types"
)

//go:generate code-gen ./ -interface=Client -struct=client
//go:generate code-gen ./ -interface=Client -tracing=pkg
//go:generate code-gen ./ -interface=Client -struct=client -tracing=pkg
//go:generate code-gen ./ -interface=Client -struct=client -timing=attr
//go:generate code-gen ./ -interface=Client -struct=client -tracing=pkg -timing=attr -o=client_gen.go

var rootConf = &flags.RootConfig{}
var rootCmd = &cobra.Command{
  Use:   "gen",
  Short: "",
  Args: cobra.NoArgs,
}

var tracingConf = &flags.TracingConfig{}
var tracingCmd = &cobra.Command{
  Use:   "tracing",
  Short: "",
  Args: cobra.NoArgs,
  Run: runCommand(func() error {
	  return tracing.Tracing(rootConf, tracingConf)
  }),
}


var timingCmd = &cobra.Command{
  Use:   "timing",
  Short: "",
  Args: cobra.NoArgs,
  Run: runCommand(func() error {
	  return timing.Timing(rootConf)
  }),
}

var constructorConf = &flags.ConstructorConfig{}
var constructorCmd = &cobra.Command{
  Use:   "constructor",
  Short: "",
  Args: cobra.NoArgs,
  Run: func(cmd *cobra.Command, args []string) {
  },
}

func runCommand(fn func() error) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args[]string) {
		err := fn()
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func main() {
	rootCmd.PersistentFlags().StringVar(&rootConf.Outdir, "outdir", "./tracing", "directory to write mocks to")
	rootCmd.PersistentFlags().StringVar(&rootConf.Outpkg, "outpkg", "tracing", "name of generated package")
	rootCmd.PersistentFlags().BoolVar(&rootConf.Inpkg, "inpkg", false, "generate a mock that goes inside the original package")
	rootCmd.PersistentFlags().StringVar(&rootConf.Dir, "dir", ".", "directory to search for interface")

	tracingCmd.Flags().StringVar(&tracingConf.InterfaceName, "name", "", "name of interface to generate wrappers for")
	tracingCmd.MarkFlagRequired("name")
	tracingCmd.MarkFlagRequired("dir")

	constructorCmd.Flags().StringVar(&constructorConf.StructName, "name", "", "name of struct to wrap constructor for")

	runCommand(func() error {
		return rootCmd.ExecuteContext(context.Background())
	})

	strct, ok := findStruct(pkg, f.structName)
	if !ok {
		log.Fatalf("Could not find struct: %s", f.structName)
	}

	var wrappers []string
	if len(f.tracingPkg) > 0 {
		if tracing.ShouldSkipInterface(iface) {
			log.Fatal("Could not find any methods taking context")
		}
		wrappers = append(wrappers, tracingWrapper)
	}
	if len(f.timingAttr) > 0 {
		if !timing.StructHasTimingAttr(strct, f.timingAttr) {
			log.Fatalf("Struct does not have specific timing attribute `%s`", f.timingAttr)
		}
		timingWrapper := timing.Gen(b, iface, f.timingAttr)
		wrappers = append(wrappers, timingWrapper)
	}
	constructor.Gen(b, iface, strct, wrappers)

	dst := os.Stdout
	if len(f.output) > 0 {
		if err := os.MkdirAll(filepath.Dir(f.output), os.ModePerm); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
		f, err := os.Create(f.output)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	if _, err := b.WriteTo(dst); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}


func findInterface(pkg *types.Package, name string) (types.Interface, bool) {
	for _, iface := range pkg.Interfaces {
		if iface.Name == name {
			return iface, true
		}
	}
	return types.Interface{}, false
}

func findStruct(pkg *types.Package, name string) (types.Struct, bool) {
	for _, strct := range pkg.Structs {
		if strct.Name == name {
			return strct, true
		}
	}
	return types.Struct{}, false
}
