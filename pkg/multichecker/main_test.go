package main

import (
	"testing"

	"github.com/boginskiy/Clicki/internal/tester"
	"github.com/boginskiy/Clicki/pkg/multichecker/staticchecker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

func testNewStaticcheck(t *testing.T, analyzer *analysis.Analyzer) {
	staticcheck := staticchecker.NewStaticcheck(analyzer)
	assert.Equal(t, len(staticcheck.Checks), 1)

	testLoadConfig(t, staticcheck)
}

func testLoadConfig(t *testing.T, staticcheck *staticchecker.Staticcheck) {
	staticcheck.LoadConfig("./staticchecker/test/staticcheck.json")
	assert.Greater(t, len(staticcheck.Checks), 0)
}

func TestMain(t *testing.T) {
	analyzer := tester.CreateAnalyzer()
	testNewStaticcheck(t, analyzer)
}
