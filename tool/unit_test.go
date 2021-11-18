package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func splitArgs(s []string) ([]string, []string) {
	for x, a := range s {
		if a == "--" {
			return s[0:x], s[x+1:]
		}
	}
	return s, []string{}
}

func TestMain(t *testing.T) {
	if os.Getenv("TEST_MAIN") == "1" {
		oldA, newA := splitArgs(os.Args)
		defer func() { os.Args = oldA }()
		os.Args = append([]string{"ir"}, newA...)
		main()
	}
}

const verbose = false

func runTest(t *testing.T, file string, exitCode int, cov string) {
	t.Run(file, func(t *testing.T) {
		entry := os.Args[0]
		cmd := exec.Command(
			entry,
			"-test.run=TestMain",
			"-test.coverprofile="+cov,
			"--",
			"-y",
			"-c",
			file,
		)
		cmd.Env = append(os.Environ(), "TEST_MAIN=1")
		o, err := cmd.CombinedOutput()
		if verbose {
			fmt.Println(string(o))
		}
		if exitCode == 0 {
			if err != nil {
				t.Errorf("process ran with err %v, want exit status 0", err)
			}
			return
		}
		e, ok := err.(*exec.ExitError)
		if !ok {
			t.Errorf("type assertion of the error failed, %v", err)
			return
		}
		if e.ExitCode() != exitCode {
			t.Errorf("process ran with exit status %v, want %d", e, exitCode)
		}
	})
}

func TestEmpty(t *testing.T) {
	for i, f := range []string{
		"../tests/md/empty001.md",
		"../tests/md/empty002.md",
		"../tests/md/empty003.md",
		"../tests/md/noentry001.md",
	} {
		runTest(t, f, exitEmpty, fmt.Sprintf("coverage%d.out", i))
	}
}

func TestFail(t *testing.T) {
	for i, f := range []string{
		"../tests/md/hello001b.md",
	} {
		runTest(t, f, exitFail, fmt.Sprintf("coverage%d.out", 100+i))
	}
}

func TestSuccess(t *testing.T) {
	for i, f := range []string{
		"../tests/md/hello001.md",
		"../tests/md/hello002.md",
		"../tests/md/docker001.md",
		"../tests/md/docker001b.md",
		"../tests/md/docker002.md",
		"../tests/md/docker003.md",
		"../tests/md/docker004.md",
		"../tests/md/attached001.md",
		//"https://raw.githubusercontent.com/umarcor/issue-runner/master/tests/md/vunit-py.md",
		//"VUnit/vunit#337",
		//"ghdl/ghdl#579",
		//"ghdl/ghdl#584",
	} {
		runTest(t, f, 0, fmt.Sprintf("coverage%d.out", 200+i))
	}
}

/*
args := []string{"-y", "-c", f[0]}
if len(f[1]) > 0 {
	args = append([]string{f[1]}, args...)
}
entry := os.Args[0]
cmd := exec.Command(entry, "-test.run=TestMain", "--")
cmd.Args = append(cmd.Args, args...)
*/
