package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/google/go-github/v28/github"
)

func getIssueBody(user, repo, id string) string {
	n, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}
	client := github.NewClient(nil)
	issue, resp, err := client.Issues.Get(context.Background(), user, repo, n)
	if err != nil {
		if _, ok := err.(*github.RateLimitError); ok {
			log.Fatal("hit GitHub API rate limit")
		}
		log.Fatal(err)
	}
	fmt.Printf("API limit: remaining %d/%d (%s)\n", resp.Rate.Remaining, resp.Rate.Limit, resp.Rate.Reset)
	return *issue.Body
}

func processArgs(args []string) (*mwes, error) {
	egs := newMWESlice(len(args))
	k := 0

	inc := func(a string) {
		egs[k].args = append(egs[k].args, a)
		if !cfgMerge {
			k++
		}
	}

	if len(args) == 1 && args[0] == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		err = egs[k].parseBody(string(data))
		if err == nil {
			inc("stdin")
		}
		return &egs, err
	}

	for x, i := range args {
		fmt.Printf("Processing arg %d/%d (%s)...\n", x+1, len(args), i)
		err := egs[k].processArg(i)
		if err != nil {
			return nil, err
		}
		inc(i)
	}
	egs = egs[0:k]
	return &egs, nil
}

func (e *mwe) processArg(arg string) error {
	fmt.Println("· processArg", arg, e)

	// Short issue format
	re := regexp.MustCompile(`(.*)/(.*)#([0-9]*)`)
	if re.MatchString(arg) {
		m := re.FindAllStringSubmatch(arg, -1)[0]
		body := getIssueBody(m[1], m[2], m[3])

		return e.parseBody(body)
	}
	// File path or URL
	if ext := filepath.Ext(arg); ext != "" {
		if ext != ".md" {
			return fmt.Errorf("WIP! non-markdown files/tarballs/zipfiles not supported as entrypoints yet")
		}
		_, err := url.ParseRequestURI(arg)
		if err == nil {
			return e.parseFromURL(arg)
		}
		return e.parseFromFile(arg)
	}
	// Otherwise, error
	return fmt.Errorf("unknown arg format %s", arg)
}

func (e *mwe) parseFromURL(arg string) error {
	fmt.Println("· parseFromURL", arg, e)
	resp, err := http.Get(arg) //nolint:gosec
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP Response Status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return e.parseBody(string(dat))
}

func (e *mwe) parseFromFile(arg string) error {
	fmt.Println("· parseFromFile", arg, e)
	info, err := os.Stat(arg)
	if !(os.IsNotExist(err) || info.IsDir()) {
		dat, err := ioutil.ReadFile(arg)
		if err != nil {
			return err
		}
		return e.parseBody(string(dat))
	}
	return fmt.Errorf("could not get file '%s'", arg)
}
