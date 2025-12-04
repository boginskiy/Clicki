package builder

import (
	"fmt"
	"go/ast"
)

type ResetStruct struct {
	Name   string
	Fields []string
	Type   []string
}

type Building struct {
	Package      string
	ResetStructs []ResetStruct
}

func NewBuilding(path string) *Building {
	return &Building{
		ResetStructs: make([]ResetStruct, 0, 10),
		Package:      path,
	}
}

func (p *Building) PutResetStruct(name string, tpy *ast.StructType) {
	tmpName := make([]string, len(tpy.Fields.List))
	// tmpType := make([]string, len(tpy.Fields.List))

	for i, field := range tpy.Fields.List {
		fmt.Println(field)
		// Name
		if len(field.Names) > 0 {
			tmpName[i] = field.Names[0].Name
		}

		// Type
		// tmpType[i] = field.Type

	}

	// Может отсюда взять инфу ?
	
	// type Field struct {
	// Doc     *CommentGroup // associated documentation; or nil
	// Names   []*Ident      // field/method/(type) parameter names; or nil
	// Type    Expr          // field/method/parameter type; or nil
	// Tag     *BasicLit     // field tag; or nil
	// Comment *CommentGroup // line comments; or nil
}

	p.ResetStructs = append(p.ResetStructs, ResetStruct{Name: name})

	fmt.Println(p.ResetStructs[0].Fields)

}

func (p *Building) Execute() {
	if p.Package == "" {
		return
	}

	// генерируем код по шаблону
	// var buf bytes.Buffer
	// err := tmpl.Execute(&buf, p)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
