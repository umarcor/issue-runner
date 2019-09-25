package main

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

func (e *mwe) parseBody(body string) error {
	md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
	tokens := md.Parse([]byte(body))

	for _, t := range tokens {
		if s := func(tok markdown.Token) (s *snippet) {
			switch tok := tok.(type) {
			case *markdown.Fence:
				s = &snippet{[]byte(tok.Content), tok.Params}
				//	case *markdown.CodeBlock:
				//		return snippet{tok.Content, "code"}
				//	case *markdown.CodeInline:
				//		return snippet{tok.Content, "code inline"}

				//	case *markdown.Inline:
				//		fmt.Println("> Inline")
				//		fmt.Println(tok.Content)
			}
			if s == nil {
				return
			}
			if ok := s.parseName(); !ok {
				s = nil
			}
			return
		}(t); s != nil {
			e.snippets = append(e.snippets, s)
		}
	}

	if links := parseLinks(body); len(links) != 0 {
		e.links = append(e.links, links...)
	} else if len(e.snippets) == 0 {
		return errEmptyBody
	}

	for _, s := range e.snippets {
		if s.name[0:2] == "i:" {
			if len(e.entry) != 0 {
				log.Fatal("entrypoint already set to", e.entry)
			}
			e.entry = strings.Split(s.name[2:], " ")
			s.name = "run"
		}
	}

	return nil
}

func (s *snippet) parseName() bool {
	for _, t := range []string{"file", "image"} {
		re := regexp.MustCompile(`.*:` + t + `:(.*)`)
		for _, c := range []string{s.name, string(s.bytes)} {
			if re.MatchString(c) {
				k := strings.TrimSpace(re.FindAllStringSubmatch(c, -1)[0][1])
				if t == "image" {
					s.name = "i:" + k
					return true
				}
				s.name = k
				return true
			}
		}
	}
	return false
}

func parseLinks(body string) []*link {
	fmt.Print("· parseLinks ")

	re := regexp.MustCompile(`\[:mwe:([^ ]*)\]\((\S*)\)`)
	if !re.MatchString(body) {
		fmt.Println("no match")
		return nil
	}
	refs := make([]*link, 0)

	s := re.FindAllStringSubmatch(body, -1)
	fmt.Println("match", len(s))
	for _, r := range s {
		if val := isValidLink(&(r[1]), &(r[2])); val {
			fmt.Printf("Found ref [%s](%s)...\n", r[1], r[2])
			refs = append(refs, &link{r[1], r[2]})
		}
	}

	return refs
}

func isValidLink(name, ref *string) bool {
	ext := filepath.Ext(*name)
	if len(ext) == 0 {
		ext = filepath.Ext(*ref)
		if len(ext) == 0 {
			fmt.Printf("WARNING! undefined extension for file <%s>, skipping", *name)
			return false
		}
		*name = filepath.Base(*ref)
	}
	if ext[1:] == "txt" {
		nname := (*name)[0 : len(*name)-len(ext)]
		if len(filepath.Ext(nname)) != 0 {
			*name = nname
		}
	}

	_, err := url.ParseRequestURI(*ref)
	return err == nil
}
