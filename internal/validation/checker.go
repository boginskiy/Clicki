package validation

import "regexp"

const (
	CheckDomain = `^(https?:)?\/\/(www\.)?[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(\/.*)?`
	CheckPath   = `^/[a-zA-Z0-9]+$`
)

// Interfaces
type Checker interface {
	CheckUpPath(path string) bool
	CheckUpURL(body string) bool
}

type Check struct {
}

func NewChecker() *Check {
	return &Check{}
}

func (t *Check) CheckUpPath(path string) bool {
	re := regexp.MustCompile(CheckPath)
	return re.MatchString(path)
}

func (t *Check) CheckUpURL(body string) bool {
	re := regexp.MustCompile(CheckDomain)
	return re.MatchString(body)
}
