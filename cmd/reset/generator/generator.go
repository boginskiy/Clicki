package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"os"
)

// Generator.
type Generator struct {
	Path    string
	Package string
	Resets  []Reset
}

func NewGenerator(path string) *Generator {
	return &Generator{
		Resets: make([]Reset, 0, 10),
		Path:   path,
	}
}

func (p *Generator) UpdateReset(name string, typeStruct *ast.StructType) {
	// Create Resets.
	tmpReset := NewReset(name, typeStruct.Fields.List)
	tmpReset.CreateDefaultFields()

	// if there is DefaultRows, add to all structs.
	if len(tmpReset.Fields) > 0 {
		p.Resets = append(p.Resets, tmpReset)
	}
}

func (p *Generator) Execute() {
	if p.Package == "" || len(p.Resets) == 0 {
		return
	}

	// Generate code about template.
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, p)
	if err != nil {
		log.Fatal("tmpl.Execute", err)
	}

	// Format code.
	bufFmt, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal("format.Source", err)
	}

	// Write code in the file.
	// basename := strings.TrimSuffix(p.Package, filepath.Ext(fname))
	err = os.WriteFile(fmt.Sprintf("%s/reset.gen.go", p.Path), bufFmt, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
