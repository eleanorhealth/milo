package main

import (
	_ "embed"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed template
var tpl string

type data struct {
	EntityName string
	EntityType string
	IDType     string
}

func main() {
	t, err := template.New("template").Funcs(template.FuncMap{
		"trimLeft": strings.TrimLeft,
	}).Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	d := data{
		EntityName: "Customer",
		EntityType: "*domain.Customer",
		IDType:     "entityid.ID",
	}

	err = t.Execute(os.Stdout, d)
	if err != nil {
		log.Fatal(err)
	}
}
