package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

var (
	version = "dev"
	commit  = "none"
	//date    = "unknown"
)

// TODO Add function to return version string

// TODO add option to fail if not interactive tty and not -y and local exec

// TODO Support providing a list of `:image: ghdl/vunit:mcode ghdl/vunit:llvm`, instead of a single image

var (
	cfgTmp   string
	cfgMerge bool
	cfgYes   bool
	cfgClean bool
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "issue-runner",
	Version: version + "-" + commit,
	Short:   "issue-runner executes MWEs from markdown files",
	Long: `Execute Minimal Working Examples (MWEs) defined in markdown files,
in the body of GitHub issues or as tarballs/zipfiles.
Site: github.com/1138-4EB/issue-runner`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := run(args, true)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var srcsCmd = &cobra.Command{
	Use:   "sources",
	Short: "extract sources but do not execute any MWE",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mwes, err := run(args, false)
		if err != nil {
			log.Fatal(err)
		}
		mwes.print()
	},
}

func init() {
	rootCmd.AddCommand(srcsCmd)
	for _, c := range []*cobra.Command{rootCmd, srcsCmd} {
		commonFlags(c.Flags())
	}
}

func commonFlags(f *flag.FlagSet) {
	f.StringVarP(&cfgTmp, "tmp", "t", "", "base directory for temporal dirs")
	f.BoolVarP(&cfgMerge, "merge", "m", false, "merge arguments in a single MWE")
	f.BoolVarP(&cfgYes, "yes", "y", false, "force response to interactive questions")
	f.BoolVarP(&cfgClean, "clean", "c", false, "remove sources after executing MWEs")
}

func run(args []string, exec bool) (*mwes, error) {
	es, err := processArgs(args)
	if err != nil {
		return nil, err
	}
	if len(*es) == 0 {
		log.Println("no MWE was found, exiting")
		return nil, err
	}
	if err := es.generate(cfgTmp); err != nil {
		return es, err
	}
	if exec {
		err := es.execute()
		if err != nil {
			return es, err
		}
	}
	if cfgClean {
		for _, e := range *es {
			fmt.Println("Removing...", e.dir)
			os.RemoveAll(e.dir)
		}
	}
	return es, nil
}
