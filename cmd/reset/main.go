package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/boginskiy/Clicki/cmd/reset/builder"
)

func main() {
	comment := "generate:reset"

	// Definition working dir.
	startDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pathsToFiles := builder.Scanner(startDir, ".go") // Get all paths to files.go .
	building := builder.NewBuilding("")              // Struct for building.

	// Iteration every files.
	for _, path := range pathsToFiles {

		// Take current package.
		currentPackage, err := builder.TakeCurrentPackage(path)
		if err != nil {
			continue
		}

		// New builder for new package.
		if building.Package != currentPackage {
			building.Execute()                             // Go to generation.
			building = builder.NewBuilding(currentPackage) // Create new builder.
		}

		// ASTree
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
		}

		// Data for Building.
		ast.Inspect(f, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.GenDecl:

				// Search all struct with comment "generate:reset" in current file.
				for _, spec := range x.Specs {
					if tspec, ok := spec.(*ast.TypeSpec); ok {

						structType, ok := tspec.Type.(*ast.StructType)
						if ok && structType != nil && strings.Contains(x.Doc.Text(), comment) {

							building.PutResetStruct(tspec.Name.String(), structType)

						}
					}
				}
			}
			return true
		})

		// TODO:убрать
		building.Execute()
		break
	}
}
