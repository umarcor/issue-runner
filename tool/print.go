package main

import (
	"fmt"
)

func (e *mwe) print() {
	fmt.Println("args:")
	for _, a := range e.args {
		fmt.Println("-", a)
	}
	fmt.Println("entrypoint:", e.entry)
	fmt.Println("directory: ", e.dir)
	fmt.Println("files:")
	for _, s := range e.snippets {
		fmt.Println("-", s.name)
	}
}

func (es *mwes) print() {
	for x, e := range *es {
		fmt.Println("\nMWE", x)
		e.print()
	}
}
