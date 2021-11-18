package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	au "github.com/logrusorgru/aurora"
	v "github.com/spf13/viper"
)

func (e *mwe) execute() error {
	loc := e.loc()
	switch len(e.entry) {
	case 0:
		return host(loc)
	case 1:
		return docker(e.entry[0], loc)
	default:
		// Multiple containers; make a copy of the sources for each
		for x, i := range e.entry {
			rdir := loc + "-" + strconv.Itoa(x)
			if err := copyDir(loc, rdir); err != nil {
				return err
			}
			if err := docker(i, rdir); err != nil {
				return err
			}
			if v.GetBool("clean") {
				fmt.Println("Removing...", rdir)
				os.RemoveAll(rdir)
			}
		}
	}
	return nil
}

func (es *mwes) execute() error {
	for x, e := range *es {
		fmt.Println("Executing...", x, e)
		e.print()
		fmt.Println(au.Cyan("|>"))
		err := e.execute()
		fmt.Println(au.Cyan("<|"))
		if err != nil {
			return err
		}
	}
	return nil
}

func askForConfirmation() bool {
	fmt.Println(":: Execute MWE on the host? [y/N]")
	var r rune
	_, err := fmt.Scanf("%c\n", &r)
	if err != nil {
		log.Fatal(err)
	}
	if r == 'y' || r == 'Y' {
		return true
	} else if r == 'n' || r == 'N' || r == '\n' || r == '\r' {
		return false
	} else {
		return askForConfirmation()
	}
}

func host(dir string) error {
	if v.GetBool("no-host") {
		return errHostExecDisabled
	}

	if !v.GetBool("yes") {
		fmt.Println(`WARNING! A MWE is about to be executed on the host. This is not recommended, since unreliable
		code can damage your system. We suggest to use an OCI container instead.`)
		if ok := askForConfirmation(); !ok {
			return nil
		}
	}

	cmd := exec.Command("sh", "./run")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	buf := &bytes.Buffer{}
	cmd.Stderr = io.MultiWriter(os.Stderr, buf)

	err := cmd.Run()

	if err != nil {
		if strings.Contains(buf.String(), "exec format error") {
			return errExecFormat
		}
		if _, ok := err.(*exec.ExitError); ok {
			return errExecFailure
		}
		return err
	}

	return nil
}
