package generator

import "html/template"

var tmpl = template.Must(template.New("reset").Parse(templateStr))

const templateStr = `
package {{.Package}}

{{range .Resets}} 

func ({{.Nickname}} *{{.Name}}) Reset() {
    if {{.Nickname}} == nil {
        return
    }

    {{range .Fields}}
        {{if and (eq .Type "string") (eq .PatternType "Var")}}
            {{.Nickname}}.{{.Name}} = ""
        {{else if eq .PatternType "Var"}}
            {{.Nickname}}.{{.Name}} = {{.Value}}

        {{else if and (eq .Type "string") (eq .PatternType "VarP")}}
            if {{.Nickname}}.{{.Name}} != nil {*{{.Nickname}}.{{.Name}} = ""}
        {{else if eq .PatternType "VarP"}}
            if {{.Nickname}}.{{.Name}} != nil {*{{.Nickname}}.{{.Name}} = {{.Value}}}
        
        {{else if eq .PatternType "StructP"}}
            if resetter, ok := any({{.Nickname}}.{{.Name}}).(interface{ Reset() }); ok && {{.Nickname}}.{{.Name}} != nil {resetter.Reset()}

        {{else}}
            {{.Pattern}}
        {{end}}
    {{end}}
}
{{end}}
`
