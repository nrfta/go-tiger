package policy

import (
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/nrfta/go-tiger/pkg/generators"
	"github.com/volatiletech/inflect"
)

func Process(rootPath string, singularResourceName string) {
	data := struct {
		Name       string
		NamePlural string
	}{
		Name:       strcase.ToCamel(singularResourceName),
		NamePlural: strcase.ToCamel(inflect.Pluralize(singularResourceName)),
	}

	gen := generators.NewGenerator(
		rootPath,
		data,
		template.FuncMap{
			"ToLowerCamel": strcase.ToLowerCamel,
		},
	)
	gen.RenderTemplateToFile("policy/policy.go.tpl", "pkg/policy/"+strcase.ToSnake(data.NamePlural)+".go")
	gen.RenderTemplateToFile("policy/policy_test.go.tpl", "pkg/policy/"+strcase.ToSnake(data.NamePlural)+"_test.go")

	gen.PrintSummary()
}
