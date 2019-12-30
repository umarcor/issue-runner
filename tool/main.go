package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/docker/docker/client"
	au "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v "github.com/spf13/viper"
)

var (
	version = "dev"
	commit  = "none"
	//date    = "unknown"
)

// TODO Add function to return version string

var (
	cfgFile               string
	errExecFailure        = fmt.Errorf("execution of the MWE failed")
	errExecFormat         = fmt.Errorf("exec format error; is there a shebang?")
	errHostExecDisabled   = fmt.Errorf("execution of MWEs on the host is disabled")
	errDockerExecDisabled = fmt.Errorf("execution of MWEs in OCI containers is disabled")
	errEmptyBody          = fmt.Errorf("no supported content found")
	errDockerConnect      = fmt.Errorf("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?")

	exitEmpty   = 1
	exitExec    = 2
	exitFormat  = 3
	exitDocker  = 4
	exitFail    = 5
	exitDefault = 6
)

func main() {
	fmt.Println(au.Sprintf(au.Cyan("[issue-runner] %s"), rootCmd.Version))
	if err := rootCmd.Execute(); err != nil {
		if client.IsErrConnectionFailed(err) ||
			client.IsErrNotFound(err) ||
			client.IsErrNotImplemented(err) {
			os.Exit(exitDocker)
		}

		switch err.Error() {
		case errEmptyBody.Error():
			os.Exit(exitEmpty)
		// TODO These two following cases might be merged in a single case statement?
		case errHostExecDisabled.Error():
			os.Exit(exitExec)
		case errDockerExecDisabled.Error():
			os.Exit(exitExec)
		case errExecFormat.Error():
			os.Exit(exitFormat)
		// TODO The following case might already be covered by IsErrConnectionFailed above
		case errDockerConnect.Error():
			os.Exit(exitDocker)
		case errExecFailure.Error():
			os.Exit(exitFail)
		default:
			os.Exit(exitDefault)
		}
	}
}

var rootCmd = &cobra.Command{
	Use:     "issue-runner",
	Version: version + "-" + commit,
	Short:   au.Sprintf(au.Cyan("issue-runner executes MWEs from markdown files")),
	Long: `Execute Minimal Working Examples (MWEs) defined in markdown files,
in the body of GitHub issues or as tarballs/zipfiles.
Site: github.com/eine/issue-runner`,
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := run(args, !v.GetBool("no-exec"))
		return err
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	f := rootCmd.Flags()
	// Helper functions to set cobra and viper at once
	flag, flagP := flagFuncs(f)

	f.StringVar(&cfgFile, "config", "", "config file (defaults are './.issue-runner[ext]', '$HOME/.issue-runner[ext]' or '/etc/issue-runner/.issue-runner[ext]')")
	flagP("tmp", "t", "", "base temporal dir")
	flagP("dir", "d", "tmp-run", "base name for temporal subdirs")
	flagP("merge", "m", false, "merge arguments in a single MWE")
	flag("no-docker", false, "disable executing MWEs in containers")
	flag("no-host", false, "disable executing MWEs on the host")
	flagP("no-exec", "x", false, "extract sources but do not execute any MWE")
	flagP("yes", "y", false, "non-interactive")
	flagP("clean", "c", false, "remove sources after executing MWEs")

	// Bind the full flag set to the configuration
	err := v.BindPFlags(f)
	if err != nil {
		log.Fatal(err)
	}
}

func flagFuncs(f *pflag.FlagSet) (flag func(k string, i interface{}, u string), flagP func(k, p string, i interface{}, u string)) {
	flag = func(k string, i interface{}, u string) {
		switch y := i.(type) {
		case bool:
			f.Bool(k, y, u)
		case int:
			f.Int(k, y, u)
		case string:
			f.String(k, y, u)
		}
		v.SetDefault(k, i)
	}
	flagP = func(k, p string, i interface{}, u string) {
		switch y := i.(type) {
		case bool:
			f.BoolP(k, p, y, u)
		case int:
			f.IntP(k, p, y, u)
		case string:
			f.StringP(k, p, y, u)
		}
		v.SetDefault(k, i)
	}
	return
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		v.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		v.AddConfigPath(".")
		v.AddConfigPath(home)
		v.AddConfigPath("/etc/issue-runner/")
		v.SetConfigName(".issue-runner")
	}

	v.SetEnvPrefix("ISSUERUNNER")
	v.AutomaticEnv()
	//v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		// Fail with invalid config format
		if _, ok := err.(v.ConfigParseError); ok {
			log.Fatal(err)
		}
	} else {
		log.Println("Using config file:", v.ConfigFileUsed())
	}

	tmp := v.GetString("tmp")
	if len(tmp) != 0 {
		info, err := os.Stat(tmp)
		if err == nil {
			if !info.IsDir() {
				log.Fatal(fmt.Errorf("'%s' exists and it is not a directory, cannot proceed", tmp))
			}
		} else if os.IsNotExist(err) {
			fmt.Println(fmt.Sprintf("MkdirAll '%s'", tmp))
			err = os.MkdirAll(tmp, 0755)
			if err != nil {
				log.Fatal(err)
			}
		} else if err != nil {
			log.Fatal(err)
		}
	}

	v.Set("indocker", false)
	if runtime.GOOS != "windows" {
		cmd := exec.Command("cat", "/proc/self/cgroup")
		o, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}
		if strings.Contains(string(o), "docker") {
			fmt.Println("it seems you are running issue-runner inside a Docker container")
			if v.GetBool("no-docker") {
				fmt.Println("but execution of sibling containers is disabled through '--no-docker'")
			} else {
				sock := "/var/run/docker.sock"
				_, err := os.Stat(sock)
				if os.IsNotExist(err) {
					log.Fatal(fmt.Errorf("'%s' does not exist", sock))
				} else if err != nil {
					log.Fatal(err)
				}

				fmt.Println("to run issues in sibling containers, ensure that 'issues:/volume' is bind")
				vol := "/volume"
				_, err = os.Stat(vol)
				if os.IsNotExist(err) {
					log.Fatal(fmt.Errorf("'%s' does not exist", vol))
				} else if err != nil {
					log.Fatal(err)
				}

				v.Set("tmp", path.Join("/volume", v.GetString("tmp")))
				v.Set("indocker", true)
			}
		}
	}

	//box = rice.MustFindBox("data")
}
