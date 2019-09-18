package main

import (
	"fmt"
	"log"
	"os/exec"
)

func (e *mwe) execute() error {
	if e.entry != "" {
		return docker(e.entry, e.dir)
	}

	fmt.Println(`WARNING! A MWE is about to be executed on the host. This is not recommended, since unreliable
code can damage your system. We suggest to use an OCI container instead.`)
	if !cfgYes {
		if ok := askForConfirmation(); !ok {
			return nil
		}
	}
	cmd := exec.Command("sh", "./run")
	cmd.Dir = e.dir
	o, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// TODO handle logs properly (tee)
	// TODO handle exit code properly
	fmt.Println(string(o))
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
