package builder

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func Scanner(startPath, ExtensionFile string) []string {
	result := make([]string, 0, 10)
	arrayFile, err := os.ReadDir(startPath)
	if err != nil || len(arrayFile) == 0 {
		return []string{}
	}

	// Iteration file.
	for _, f := range arrayFile {
		// Check folder.
		if f.IsDir() {
			result = append(result, Scanner(filepath.Join(startPath, f.Name()), ExtensionFile)...)
			// File is not '.go' .
		} else if filepath.Ext(f.Name()) != ExtensionFile {
			continue
			// File is '.go' .
		} else {
			result = append(result, filepath.Join(startPath, f.Name()))
		}
	}
	return result
}

func TakeCurrentPackage(path string) (string, error) {
	tmpPath := strings.Split(path, "/")

	if len(tmpPath) < 2 {
		return "", errors.New("file out of package. Will be ignore")
	}

	tmpPath = tmpPath[:len(tmpPath)-1]
	return strings.Join(tmpPath, "/"), nil
}
