package main

import (
	_ "embed"
	"flag"
	"log"
	"os"
	"strings"
	"text/template"
)

//go:embed template
var tpl string

type data struct {
	EntityName        string
	EntityType        string
	IDType            string
	NotFoundErrorType string
}

func main() {
	var entityName, entityType, idType, notFoundErrorType string

	flag.StringVar(&entityName, "entityName", "", "domain entity name (e.g., Customer)")
	flag.StringVar(&entityType, "entityType", "", "domain entity type (e.g., *domain.Customer)")
	flag.StringVar(&idType, "idType", "", "domain ID type (e.g., entityid.ID)")
	flag.StringVar(&notFoundErrorType, "notFoundErrorType", "", "entity not found error type (e.g., domain.ErrNotFound)")

	flag.Parse()

	if len(entityName) == 0 || len(entityType) == 0 || len(idType) == 0 || len(notFoundErrorType) == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	t, err := template.New("template").Funcs(template.FuncMap{
		"trimLeft": strings.TrimLeft,
	}).Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	d := data{
		EntityName:        entityName,
		EntityType:        entityType,
		IDType:            idType,
		NotFoundErrorType: notFoundErrorType,
	}

	err = t.Execute(os.Stdout, d)
	if err != nil {
		log.Fatal(err)
	}
}
