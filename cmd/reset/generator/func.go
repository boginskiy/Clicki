package generator

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func scannerTestFile(file, substr string) bool {
	return strings.Contains(file, substr)
}

func Scanner(startPath, extensionFile, exclude string) []string {
	result := make([]string, 0, 10)
	arrayFile, err := os.ReadDir(startPath)
	if err != nil || len(arrayFile) == 0 {
		return []string{}
	}

	// Iteration file.
	for _, f := range arrayFile {
		// Check folder.
		if f.IsDir() {
			result = append(result, Scanner(filepath.Join(startPath, f.Name()), extensionFile, exclude)...)
			// File is not '.go' .
		} else if filepath.Ext(f.Name()) != extensionFile {
			continue
			// File is '.go' .

		} else {
			// if file.go, adding it.
			if !scannerTestFile(f.Name(), exclude) {
				result = append(result, filepath.Join(startPath, f.Name()))
			}
			// if file_test.go, passing it.
		}
	}
	return result
}

func TakePath(path string) (string, error) {
	tmpPath := strings.Split(path, "/")

	if len(tmpPath) < 2 {
		return "", errors.New("file out of package. Will be ignore")
	}

	tmpPath = tmpPath[:len(tmpPath)-1]
	return strings.Join(tmpPath, "/"), nil
}

func Pprint(fieldType ast.Expr) {
	buf := new(strings.Builder)
	ast.Fprint(buf, token.NewFileSet(), fieldType, nil)
	fmt.Println(buf.String())
}
