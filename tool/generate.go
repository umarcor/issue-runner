package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func (e *mwe) generate(tmp string) error {
	// FIXME use 'tmp' if not empty, and update e.dir accordingly
	loc := e.dir
	info, err := os.Stat(loc)

	if os.IsNotExist(err) {
		err := os.MkdirAll(loc, 0755)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !info.IsDir() {
		return fmt.Errorf("%s exists and it is not a directory, cannot proceed", loc)
	}

	for _, s := range e.snippets {
		p := path.Join(loc, s.name)
		if _, err := os.Stat(p); err == nil {
			fmt.Println("WARNING! file", p, "exists, overwriting")
		}
		err := ioutil.WriteFile(p, s.bytes, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func (es *mwes) generate(tmp string) error {
	for x, e := range *es {
		fmt.Println("Generating...", x, e)
		if err := e.generate(tmp); err != nil {
			return err
		}
	}
	return nil
}
