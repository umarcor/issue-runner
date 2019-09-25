package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func (e *mwe) execute() error {
	switch len(e.entry) {
	case 0:
		return host(e.dir)
	case 1:
		return docker(e.entry[0], e.dir)
	default:
		for x, i := range e.entry {
			rdir := e.dir + "-" + strconv.Itoa(x)
			if err := copyDir(e.dir, rdir); err != nil {
				return err
			}
			if err := docker(i, rdir); err != nil {
				return err
			}
			if cfgClean {
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
		err := e.execute()
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
	if cfgNo {
		return errHostExecDisabled
	}

	if !cfgYes {
		fmt.Println(`WARNING! A MWE is about to be executed on the host. This is not recommended, since unreliable
		code can damage your system. We suggest to use an OCI container instead.`)
		if ok := askForConfirmation(); !ok {
			return nil
		}
	}

	cmd := exec.Command("sh", "./run")
	cmd.Dir = dir
	o, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// TODO handle logs properly (tee)
	// TODO handle exit code properly
	fmt.Println(string(o))
	return nil
}
