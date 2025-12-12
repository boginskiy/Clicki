package staticchecker

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// Staticcheck is checker for Static analize code.
type Staticcheck struct {
	Checks []*analysis.Analyzer
}

func NewStaticcheck(checks ...*analysis.Analyzer) *Staticcheck {
	return &Staticcheck{
		Checks: checks,
	}
}

func (s *Staticcheck) updateChecks(set map[string]struct{}) []*analysis.Analyzer {
	var checks []*analysis.Analyzer

	// Add all Analyzers.
	staticcheck.Analyzers = append(staticcheck.Analyzers, quickfix.Analyzers...)
	staticcheck.Analyzers = append(staticcheck.Analyzers, simple.Analyzers...)
	staticcheck.Analyzers = append(staticcheck.Analyzers, stylecheck.Analyzers...)

	for _, v := range staticcheck.Analyzers {
		// Processing usual code check. Ex. "SA6003".
		if _, ok := set[v.Analyzer.Name]; ok {
			checks = append(checks, v.Analyzer)

		} else {
			// Processing special code Ex. "SA1...".
			rName := []rune(v.Analyzer.Name)
			if len(rName) >= 3 {
				copy(rName[len(rName)-3:], []rune("..."))
			}

			if _, ok := set[string(rName)]; ok {
				checks = append(checks, v.Analyzer)
			}
		}
	}
	return checks
}

// LoadConfig for loading of config params from file.
func (s *Staticcheck) LoadConfig(pathToFile string) {
	// Read config.
	config := &Config{}
	data := readFile(pathToFile)
	deserialization(data, config)

	// Create Set.
	setCheck := make(map[string]struct{})
	for _, v := range config.Staticcheck {
		setCheck[v] = struct{}{}
	}

	// Update Checks.
	s.Checks = append(s.Checks, s.updateChecks(setCheck)...)
}

func (s *Staticcheck) Start() {
	multichecker.Main(s.Checks...)
}

func (s *Staticcheck) SetChecks(checks ...*analysis.Analyzer) []*analysis.Analyzer {
	s.Checks = append(s.Checks, checks...)
	return s.Checks
}

func (s *Staticcheck) DeleteCheck(name string) []*analysis.Analyzer {
	return nil
}
