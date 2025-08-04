package pkg

import "regexp"

const (
	// Some regular expressions
	CheckDomain = `^(https?:)?\/\/(www\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(\/.*)?`
	CheckPath   = `^/[a-zA-Z0-9]+$`
)

type Tools struct {
}

func (t *Tools) CheckUpPath(path string) bool {
	re := regexp.MustCompile(CheckPath)
	return re.MatchString(path)
}

func (t *Tools) CheckUpBody(body string) bool {
	re := regexp.MustCompile(CheckDomain)
	return re.MatchString(body)
}
