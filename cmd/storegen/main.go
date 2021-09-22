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

//go:embed template-tests
var tplTests string

type data struct {
	EntityName        string
	EntityType        string
	IDType            string
	NotFoundErrorType string
}

func main() {
	var entityName, entityType, idType, notFoundErrorType string
	var tests bool

	flag.StringVar(&entityName, "entityName", "", "domain entity name (e.g., Customer)")
	flag.StringVar(&entityType, "entityType", "", "domain entity type (e.g., *domain.Customer)")
	flag.StringVar(&idType, "idType", "", "domain ID type (e.g., entityid.ID)")
	flag.StringVar(&notFoundErrorType, "notFoundErrorType", "", "entity not found error type (e.g., domain.ErrNotFound)")
	flag.BoolVar(&tests, "tests", false, "generate code for tests (assumes MockMiloStorer as the type for the mock milo.Storer)")

	flag.Parse()

	if len(entityName) == 0 || len(entityType) == 0 || len(idType) == 0 || len(notFoundErrorType) == 0 {
		flag.PrintDefaults()
		os.Exit(0)
	}

	funcs := template.FuncMap{
		"trimLeft": strings.TrimLeft,
	}

	t, err := template.New("template").Funcs(funcs).Parse(tpl)
	if err != nil {
		log.Fatal(err)
	}

	tTests, err := template.New("template-test").Funcs(funcs).Parse(tplTests)
	if err != nil {
		log.Fatal(err)
	}

	d := data{
		EntityName:        entityName,
		EntityType:        entityType,
		IDType:            idType,
		NotFoundErrorType: notFoundErrorType,
	}

	if tests {
		err = tTests.Execute(os.Stdout, d)
	} else {
		err = t.Execute(os.Stdout, d)
	}

	if err != nil {
		log.Fatal(err)
	}
}
