package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChecker(t *testing.T) {
	checker := NewChecker()

	testCheckUpPath(t, checker)
	testCheckUpURL(t, checker)
}

func testCheckUpPath(t *testing.T, checker *Check) {
	tests := []struct {
		name   string
		path   string
		result bool
	}{
		{"negative test func of CheckUpPath ", "/*******", false},
		{"positive test func of CheckUpPath ", "/asd4Gaxk", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := checker.CheckUpPath(tt.path)
			assert.Equal(t, res, tt.result)
		})
	}
}

func testCheckUpURL(t *testing.T, checker *Check) {
	tests := []struct {
		name   string
		path   string
		result bool
	}{
		{"negative test func of CheckUpURL ", "https://github.", false},
		{"positive test func of CheckUpURL ", "https://github.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := checker.CheckUpURL(tt.path)
			assert.Equal(t, res, tt.result)
		})
	}
}
