package staticchecker

import (
	"testing"

	"github.com/boginskiy/Clicki/internal/tester/tfunc"
	"github.com/stretchr/testify/assert"
)

func TestStaticcheck(t *testing.T) {
	staticcheck := NewStaticcheck()

	testUpdateChecks(t, staticcheck)
	testLoadConfig(t, staticcheck)
	testSetChecks(t, staticcheck)
	testDeleteCheck(t, staticcheck)
}

func testUpdateChecks(t *testing.T, staticcheck *Staticcheck) {
	tests := []struct {
		name     string
		set      map[string]struct{}
		expected int
	}{
		{"test adding usual code", map[string]struct{}{"SA1000": struct{}{}, "SA6000": struct{}{}}, 2},
		{"test adding code with ...", map[string]struct{}{"SA6...": struct{}{}}, 6},
		{"test adding bad code", map[string]struct{}{"AA1111": struct{}{}}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := staticcheck.updateChecks(tt.set)
			assert.Equal(t, tt.expected, len(res))
		})
	}
}

func testLoadConfig(t *testing.T, staticcheck *Staticcheck) {
	staticcheck.LoadConfig("test/staticcheck.json")
	assert.NotEmpty(t, staticcheck.Checks)
}

func testSetChecks(t *testing.T, staticcheck *Staticcheck) {
	analyzer := tfunc.CreateAnalyzer()
	staticcheck.SetChecks(analyzer)
	assert.NotEmpty(t, staticcheck.Checks)
}

func testDeleteCheck(t *testing.T, staticcheck *Staticcheck) {
	res := staticcheck.DeleteCheck("name")
	assert.Nil(t, res)
}
