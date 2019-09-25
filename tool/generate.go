package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/mholt/archiver"
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

	for _, l := range e.links {
		resp, err := http.Get((*l)[1])
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		fname := (*l)[0]
		_, err = archiver.ByExtension(fname)
		if err != nil {
			if err = bodyToDisk(resp.Body, path.Join(loc, fname)); err != nil {
				return err
			}
			continue
		}
		if err = unarchive(resp.Body, loc, fname); err != nil {
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

func bodyToDisk(src io.Reader, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)
	return err
}

func unarchive(src io.Reader, loc, name string) error {
	fpath := path.Join(loc, name)
	if err := bodyToDisk(src, fpath); err != nil {
		return err
	}
	return archiver.Unarchive(fpath, loc)
}
