package customanalyzer

import (
	"go/parser"
	"go/token"
	"log"
	"testing"

	"golang.org/x/tools/go/analysis"
)

func TestErrUsingOsExit(t *testing.T) {
	// var src = `
	// 	package main

	// 	import "fmt"

	// 	func main() {
	// 		fmt.Println("Hello, world!")
	// 		os.Exit(0)
	// 	}
	// `

	fileSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fileSet, "ErrUsingOsExit.go", nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("error of parsing: %v", err)
	}

	pass := &analysis.Pass{}
	pass.Files = append(pass.Files, astFile)
	ErrUsingOsExit.Run(pass)
}

// TODO: что за х. ?
// Work
// >> &{Doc:<nil> Package:1 Name:customanalyzer Decls:[0xc000038c00 0xc000038cc0 0xc000109110] FileStart:1 FileEnd:1198 Scope:scope 0xc00002c780 {
//         var ErrUsingOsExit
//         func run
// }
//  Imports:[0xc000108b70 0xc000108ba0 0xc000108bd0] Unresolved:[analysis analysis error fmt ast ast bool ast true false ast true false ast ast ast false true nil nil] Comments:[0xc00000e180 0xc00000e270 0xc00000e318 0xc00000e378 0xc00000e3d8 0xc00000e438 0xc00000e498 0xc00000e5d0] GoVersion:}
// <<

// Doesn't work
// >> &{Doc:<nil> Package:4 Name:main Decls:[0xc0000a40c0 0xc0000bc1b0] FileStart:1 FileEnd:101 Scope:scope 0xc0000a6180 {
//         func main
// }
//  Imports:[0xc0000bc120] Unresolved:[fmt os] Comments:[] GoVersion:}
// <<
