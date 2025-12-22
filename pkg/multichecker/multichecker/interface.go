package multichecker

import "golang.org/x/tools/go/analysis"

type Starter interface {
	Start()
}

type Setter interface {
	SetChecks(checks ...*analysis.Analyzer) []*analysis.Analyzer
}

type Deleter interface {
	DeleteCheck(name string) []*analysis.Analyzer
}

type Louder interface {
	LoadConfig(pathToFile string)
}

type Multichecker interface {
	Starter
}
