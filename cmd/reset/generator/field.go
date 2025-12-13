package generator

import (
	"fmt"
	"go/ast"
)

var DefaultValues = map[string]string{
	"int":     "0",
	"float64": "0.0",
	"float32": "0.0",
	"string":  "",
	"bool":    "false",
}

type Field struct {
	Name        string
	Value       string
	Type        string
	Pattern     string
	PatternType string
	Nickname    string
	tp          ast.Expr
}

func NewField(f *ast.Field, nickname string) *Field {
	return &Field{
		Name:     f.Names[0].Name,
		Nickname: nickname,
		tp:       f.Type,
	}
}

func (d *Field) generatVar() string {
	if x, ok := d.tp.(*ast.Ident); ok {
		if v, ok := DefaultValues[x.Name]; ok {

			d.PatternType = "Var"
			d.Type = x.Name
			d.Value = v

			d.Pattern = fmt.Sprintf(`%s.%s = %s`, d.Nickname, d.Name, DefaultValues[x.Name])
			return d.Pattern
		}
	}
	return ""
}

func (d *Field) generatVarP() string {
	if x, ok := d.tp.(*ast.StarExpr); ok {
		if p, ok := x.X.(*ast.Ident); ok {
			if v, ok := DefaultValues[p.Name]; ok {

				d.PatternType = "VarP"
				d.Type = p.Name
				d.Value = v

				d.Pattern = fmt.Sprintf(
					`if %s.%s != nil {
					*%s.%s = %s}`, d.Nickname, d.Name, d.Nickname, d.Type, d.Value)

				return d.Pattern
			}
		}
	}
	return ""
}

func (d *Field) generatArray() string {
	if x, ok := d.tp.(*ast.ArrayType); ok {
		if p, ok := x.Elt.(*ast.Ident); ok {
			if v, ok := DefaultValues[p.Name]; ok {

				d.PatternType = "Arr"
				d.Type = p.Name
				d.Value = v

				d.Pattern = fmt.Sprintf("%s.%s = %s.%s[:0]", d.Nickname, d.Name, d.Nickname, d.Name)
				return d.Pattern
			}
		}
	}
	return ""
}

func (d *Field) generatMap() string {
	if _, ok := d.tp.(*ast.MapType); ok {
		d.Pattern = fmt.Sprintf("clear(%s.%s)", d.Nickname, d.Name)
		return d.Pattern
	}
	return ""
}

func (d *Field) generatStructP() string {
	if x, ok := d.tp.(*ast.StarExpr); ok {
		if y, ok2 := x.X.(*ast.Ident); ok2 && y.Obj != nil {

			d.PatternType = "StructP"

			d.Pattern = fmt.Sprintf(
				`if resetter, ok := %s.%s.(interface{ Reset() }); ok && %s.%s != nil {
				resetter.Reset()}`,
				d.Nickname, d.Name, d.Nickname, d.Name)

			return d.Pattern
		}
	}
	return ""
}

// Generate
func (d *Field) Generate() string {
	if s := d.generatVar(); s != "" {
		return s
	}

	if s := d.generatVarP(); s != "" {
		return s
	}

	if s := d.generatArray(); s != "" {
		return s
	}

	if s := d.generatMap(); s != "" {
		return s
	}

	if s := d.generatStructP(); s != "" {
		return s
	}

	return ""
}
