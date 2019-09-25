package main

import (
	"io"
	"io/ioutil"
	"os"
	"path"
)

func copy(src, dst, name string, isdir bool) error {
	srcfp := path.Join(src, name)
	dstfp := path.Join(dst, name)
	if isdir {
		return copyDir(srcfp, dstfp)
	}
	return copyFile(srcfp, dstfp)
}

// copy a directory recursively
func copyDir(src string, dst string) error {
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(dst, srcinfo.Mode()); err != nil {
		return err
	}
	fds, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}
	for _, fd := range fds {
		err = copy(src, dst, fd.Name(), fd.IsDir())
		if err != nil {
			return err
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcfd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcfd.Close()
	dstfd, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstfd.Close()
	_, err = io.Copy(dstfd, srcfd)
	if err != nil {
		return err
	}
	srcinfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcinfo.Mode())
}
