package main

import (
	"strconv"

	v "github.com/spf13/viper"
)

type mwes []*mwe

type mwe struct {
	args     []string
	entry    []string
	dir      string
	links    []*link
	snippets []*snippet
}

type link [2]string

type snippet struct {
	bytes []byte
	name  string
}

func newMWE(d string) *mwe {
	return &mwe{
		make([]string, 0),
		make([]string, 0),
		d,
		make([]*link, 0),
		make([]*snippet, 0),
	}
}

func newMWESlice(n int) (es mwes) {
	if v.GetBool("merge") || n == 1 {
		es = make([]*mwe, 1)
		es[0] = newMWE(v.GetString("dir"))
	} else {
		es = make([]*mwe, n)
		for x := range es {
			es[x] = newMWE(v.GetString("dir") + "-" + strconv.Itoa(x))
		}
	}
	return es
}
