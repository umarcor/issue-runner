package main

import (
	"fmt"
	"log"
	"os"

	au "github.com/logrusorgru/aurora"
	v "github.com/spf13/viper"
)

func run(args []string, exec bool) (*mwes, error) {
	fmt.Println("Run", args, "exec:", exec)

	fmt.Println(au.Cyan("> Process"))
	es, err := processArgs(args)
	if err != nil {
		return nil, err
	}
	if len(*es) == 0 {
		log.Println("no MWE was found, exiting")
		return nil, err
	}

	fmt.Println(au.Cyan("> Generate"))
	if err := es.generate(); err != nil {
		return es, err
	}

	if exec {
		fmt.Println(au.Cyan("> Execute"))
		err := es.execute()
		if err != nil {
			return es, err
		}
	}

	if v.GetBool("clean") {
		for _, e := range *es {
			loc := e.loc()
			fmt.Println("Removing...", loc)
			os.RemoveAll(loc)
		}
	}
	return es, nil
}
