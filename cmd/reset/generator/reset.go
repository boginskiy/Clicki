package generator

import (
	"go/ast"
	"strings"
	"unicode"
)

// Reset
type Reset struct {
	Name       string
	Nickname   string
	Fields     []Field
	fieldsList []*ast.Field
}

func NewReset(name string, fieldsList []*ast.Field) Reset {
	tmp := Reset{
		Name:       name,
		fieldsList: fieldsList,
		Fields:     make([]Field, 0, len(fieldsList))}

	tmp.assemblNickname(name)
	return tmp
}

func (r *Reset) CreateDefaultFields() {
	for _, field := range r.fieldsList {
		// Add field & default value.
		if len(field.Names) > 0 {

			// Pprint(field.Type)
			tmpField := NewField(field, r.Nickname)

			if tmpField.Generate() != "" {
				r.addDefaultField(*tmpField)
			}
			// fmt.Println(row)
		}
	}
}

func (r *Reset) addDefaultField(field Field) {
	r.Fields = append(r.Fields, field)
}

func (r *Reset) assemblNickname(name string) {
	res := make([]string, 0, 2)

	for i, j := 0, 0; i < len(name) && j < cap(res); i++ {
		if ok := unicode.IsUpper(rune(name[i])); ok {
			res = append(res, strings.ToLower(string(name[i])))
			j++
		}
	}
	r.Nickname = strings.Join(res, "")
}
