package builder

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"os"
	"strings"
	"unicode"
)

// ResetStruct
type ResetStruct struct {
	Name        string
	Nickname    string
	DefaultVals []DefaultVal
}

func NewResetStruct(name string, cnt int) ResetStruct {
	tmp := ResetStruct{Name: name, DefaultVals: make([]DefaultVal, 0, cnt)}
	tmp.assemblNickname(name)
	return tmp
}

func (r *ResetStruct) AddDefaultVal(value DefaultVal) {
	r.DefaultVals = append(r.DefaultVals, value)
}

func (r *ResetStruct) assemblNickname(name string) {
	res := make([]string, 0, 2)

	for i, j := 0, 0; i < len(name) && j < cap(res); i++ {
		if ok := unicode.IsUpper(rune(name[i])); ok {
			res = append(res, strings.ToLower(string(name[i])))
			j++
		}
	}
	r.Nickname = strings.Join(res, "")
}

// Building.
type Building struct {
	Path         string
	Package      string
	ResetStructs []ResetStruct
}

func NewBuilding(path string) *Building {
	return &Building{
		ResetStructs: make([]ResetStruct, 0, 10),
		Path:         path,
	}
}

func (p *Building) PutResetStruct(name string, typeStruct *ast.StructType) {
	resetSt := NewResetStruct(name, len(typeStruct.Fields.List)) // Create ResetStruct.

	for _, field := range typeStruct.Fields.List {
		// Add field & default value.
		if len(field.Names) > 0 {

			Pprint(field.Type)
			defaultVal := NewDefaultVal(field)

			resetSt.AddDefaultVal(*defaultVal)
		}
	}
	// if there is DefaultVals, add to all structs.
	if len(resetSt.DefaultVals) > 0 {
		p.ResetStructs = append(p.ResetStructs, resetSt)
	}
}

func (p *Building) Execute() {
	if p.Package == "" || len(p.ResetStructs) == 0 {
		return
	}

	// генерируем код по шаблону
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, p)
	if err != nil {
		log.Fatal("tmpl.Execute", err)
	}

	// форматируем код
	bufFmt, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal("format.Source", err)
	}

	// записываем сгенерированный код в файл
	// basename := strings.TrimSuffix(p.Package, filepath.Ext(fname))
	err = os.WriteFile(fmt.Sprintf("%s/test_test.go", p.Path), bufFmt, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
