package main

import (
	"strconv"
)

type mwes []*mwe

type mwe struct {
	args     []string
	entry    string
	dir      string
	snippets []*snippet
}

type snippet struct {
	bytes []byte
	name  string
}

func newMWE(d string) *mwe {
	return &mwe{
		make([]string, 0),
		"",
		d,
		make([]*snippet, 0),
	}
}

func newMWESlice(n int) (es mwes) {
	if cfgMerge {
		es = make([]*mwe, 1)
		es[0] = newMWE("tmp-run")
	} else {
		es = make([]*mwe, n)
		for x := range es {
			es[x] = newMWE("tmp-run-" + strconv.Itoa(x))
		}
	}
	return es
}
