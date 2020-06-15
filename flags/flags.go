package flags

type RootConfig struct {
	Outdir string
	Outpkg string
	Inpkg  bool
	Dir    string
}

type TracingConfig struct {
	InterfaceName string
}

type TimingConfig struct {
	InterfaceName string
}

type ConstructorConfig struct {
	StructName string
}

func Bind(rootCmd, tracingCmd, timingCmd, constructorCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&rootConf.Outdir, "outdir", "./tracing", "directory to write mocks to")
	rootCmd.PersistentFlags().StringVar(&rootConf.Outpkg, "outpkg", "tracing", "name of generated package")
	rootCmd.PersistentFlags().BoolVar(&rootConf.Inpkg, "inpkg", false, "generate a mock that goes inside the original package")
	rootCmd.PersistentFlags().StringVar(&rootConf.Dir, "dir", ".", "directory to search for interface")

	tracingCmd.Flags().StringVar(&tracingConf.InterfaceName, "name", "", "name of interface to generate wrappers for")
	tracingCmd.MarkFlagRequired("name")
	tracingCmd.MarkFlagRequired("dir")

	constructorCmd.Flags().StringVar(&constructorConf.StructName, "name", "", "name of struct to wrap constructor for")
      }
