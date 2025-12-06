package builder

import (
	"fmt"
	"go/ast"
)

type DefaultVal struct {
	Field     string
	Default   string
	Structure string
}

func NewDefaultVal(field *ast.Field) *DefaultVal {
	// field.Name
	tmp := &DefaultVal{
		Field: field.Names[0].Name,
	}
	// field.Type
	tmp.definDefaultVal(field.Type)
	return tmp
}

func (d *DefaultVal) definIdent(x *ast.Ident) {
	deft := ""
	stre := ""

	switch x.Name {
	case "int":
		deft, stre = "0", fmt.Sprintln(d.Field, "=", "0")
	case "float64":
		deft, stre = "0.0", fmt.Sprintln(d.Field, "=", "0.0")
	case "string":
		deft, stre = "", fmt.Sprintln(d.Field, "=", "")
	case "bool":
		deft, stre = "false", fmt.Sprintln(d.Field, "=", "false")
	}

}

// definDefaultVal
func (d *DefaultVal) definDefaultVal(fieldType ast.Expr) string {
	switch x := fieldType.(type) {
	case *ast.Ident:
		d.definIdent(x)
	case *ast.ArrayType:
		fmt.Println(x)
	}
	return ""
}
