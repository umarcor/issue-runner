package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"gitlab.com/golang-commonmark/markdown"
)

var errEmptyBody = fmt.Errorf("no supported content found")

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

	if len(e.snippets) < 1 {
		return errEmptyBody
	}

	for _, s := range e.snippets {
		if s.name[0:2] == "i:" {
			if e.entry != "" {
				log.Fatal("entrypoint already set to", e.entry)
			}
			e.entry = s.name[2:]
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
