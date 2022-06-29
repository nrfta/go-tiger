package factories

import (
	"{{.ModuleName}}/pkg/models"

	"github.com/kolach/go-factory"
)

// {{.Name}} factory
var {{.Name}} = factory.NewFactory(
	models.{{.Name}}{},

	// TODO
	{{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// factory.Use({{$field.Type}}).For("{{$field.Name}}"),{{end}}
	{{ end }}
)
