package customanalyzer

import (
	"go/parser"
	"go/token"
	"log"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestErrUsingOsExit(t *testing.T) {
	var src = `
		package main

		import "fmt"

		func main() {
			fmt.Println("Hello, world!")
			os.Exit(0)
		}
	`

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, "", src, parser.ParseComments)
	if err != nil {
		log.Fatalf("error of parsing: %v", err)
	}

	pass := &analysis.Pass{}
	pass.Files = append(pass.Files, astFile)
	ErrUsingOsExit.Run(pass)
}
