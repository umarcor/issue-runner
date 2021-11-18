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

func TestEmpty(t *testing.T) {
	for _, f := range []string{
		"001",
		"002",
		"003",
	} {
		t.Run(f, func(t *testing.T) {
			entry := os.Args[0]
			cmd := exec.Command(entry, "-test.run=TestMain", "--", "-y", "-c", "../tests/md/empty"+f+".md")
			cmd.Env = append(os.Environ(), "TEST_MAIN=1")
			o, err := cmd.CombinedOutput()
			fmt.Println(string(o))
			e, ok := err.(*exec.ExitError)
			if !ok || e.ExitCode() != 1 {
				t.Errorf("process ran with err %v, want exit status 1", err)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	for _, f := range [][2]string{
		{"../tests/md/hello001.md", ""},
		{"../tests/md/hello001b.md", ""},
		{"../tests/md/hello002.md", ""},
		{"../tests/md/hello003.md", "-n"},
		{"../tests/md/hello004.md", "-n"},
		{"../tests/md/hello005.md", ""},
		{"../tests/md/hello006.md", ""},
		{"../tests/md/attached001.md", ""},
		{"https://raw.githubusercontent.com/eine/issue-runner/master/tests/md/vunit-py.md", ""},
		{"VUnit/vunit#337", ""},
		{"ghdl/ghdl#579", ""},
		{"ghdl/ghdl#584", ""},
	} {
		t.Run(f[0], func(t *testing.T) {
			args := []string{"-y", "-c", f[0]}
			if len(f[1]) > 0 {
				args = append([]string{f[1]}, args...)
			}
			entry := os.Args[0]
			cmd := exec.Command(entry, "-test.run=TestMain", "--")
			cmd.Args = append(cmd.Args, args...)
			cmd.Env = append(os.Environ(), "TEST_MAIN=1")
			o, err := cmd.CombinedOutput()
			fmt.Println(string(o))
			if err != nil {
				t.Errorf("process ran with err %v, want exit status 0", err)
			}
		})
	}
}
