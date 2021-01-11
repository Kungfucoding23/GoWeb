package main

import (
	"os"
	"text/template"
)

// Persona ,..
type Persona struct {
	Nombre string
	Edad   int
}

const tp = `
{{range .}}
	{{if .Edad}}
		Nombre: {{.Nombre}} - Edad: {{.Edad}} - Correcto
		{{else if .Nombre}}
		Nombre: {{.Nombre}} - Edad: {{.Edad}} - sin edad
		{{else}}
		Nombre: {{.Nombre}} - Edad: {{.Edad}} - sin nombre ni edad
	{{end}}
{{end}}`

func main() {
	persona := []Persona{
		{"Alejandro", 24},
		{"Pedro", 25},
		{"Maria", 36},
		{"Robert", 0},
		{"", 0},
	}
	t := template.New("persona")
	t, err := t.Parse(tp)
	if err != nil {
		panic(err)
	}
	err = t.Execute(os.Stdout, persona)
	if err != nil {
		panic(err)
	}
}
