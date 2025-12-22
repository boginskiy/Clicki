package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/boginskiy/Clicki/cmd/reset/generator"
)

func main() {
	// Definition working dir.
	startDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pathsToFiles := generator.Scanner(startDir, ".go", "_test.go") // Get all paths to files.go .
	gen := generator.NewGenerator("")                              // Struct for codeGen.

	// Params
	comment := "generate:reset"

	// Iteration every files.
	for _, path := range pathsToFiles {

		// Take path to folder.
		pathToFolder, err := generator.TakePath(path)
		if err != nil {
			continue
		}

		// New generator for new package.
		if gen.Path != pathToFolder {
			gen.Execute()                              // Go to generation.
			gen = generator.NewGenerator(pathToFolder) // Create new generator.
		}

		// ASTree
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Fatal(err)
		}

		// Data for CodeGen.
		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {

			case *ast.File:
				// Add package in Struct.
				gen.Package = x.Name.Name

			case *ast.GenDecl:

				// Search all struct with comment "generate:reset" in current file.
				for _, spec := range x.Specs {
					if tspec, ok := spec.(*ast.TypeSpec); ok {

						structType, ok := tspec.Type.(*ast.StructType)
						if ok && structType != nil && strings.Contains(x.Doc.Text(), comment) {

							gen.UpdateReset(tspec.Name.String(), structType)

						}
					}
				}
			}
			return true
		})
		gen.Execute()
		break
	}
}
