package multichecker

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

// Multicheck is checker for good check static code.
type Multicheck struct {
	Checks []*analysis.Analyzer
}

func NewMulticheck(checks ...*analysis.Analyzer) *Multicheck {
	return &Multicheck{
		Checks: checks,
	}
}

// Start Multicheck.
func (m *Multicheck) Start() {
	multichecker.Main(m.Checks...)
}

// AddChecks.
func (m *Multicheck) SetChecks(checks ...*analysis.Analyzer) []*analysis.Analyzer {
	m.Checks = append(m.Checks, checks...)
	return m.Checks
}

// GetChecks.
func (m *Multicheck) GetChecks() []*analysis.Analyzer {
	return m.Checks
}

// DeleteCheck.
func (m *Multicheck) DeleteCheck(name string) []*analysis.Analyzer {
	for i, ch := range m.Checks {
		if ch.Name == name {
			m.Checks = append(m.Checks[:i], m.Checks[i+1:]...)
		}
	}
	return m.Checks
}
