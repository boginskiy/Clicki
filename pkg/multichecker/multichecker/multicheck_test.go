package multichecker

import (
	"testing"

	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/stretchr/testify/assert"
)

func TestMulticheck(t *testing.T) {
	multicheck := NewMulticheck()

	testSetChecks(t, multicheck)
	testDeleteCheck(t, multicheck)
}

func testSetChecks(t *testing.T, multicheck *Multicheck) {
	analyzer := tester.CreateAnalyzer()
	multicheck.SetChecks(analyzer)
	assert.NotEmpty(t, multicheck.Checks)
}

func testDeleteCheck(t *testing.T, multicheck *Multicheck) {
	res := multicheck.DeleteCheck("Test")
	assert.Empty(t, res)
}
