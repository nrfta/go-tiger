package scaffolding

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcMap template.FuncMap = template.FuncMap{
	"ToLowerCamel":    strcase.ToLowerCamel,
	"GraphqlFields":   GraphqlFields,
	"GraphqlField":    GraphqlField,
	"isReadOnlyField": isReadOnlyField,
}

var mapGQLTypes = map[string]string{
	"string":            "String!",
	"int":               "Int!",
	"bool":              "Boolean!",
	"null.String":       "String",
	"time.Time":         "DateTime!",
	"null.Time":         "DateTime",
	"types.StringArray": "[String!]!",
}

func GraphqlField(field fieldDef, mustBeOptional bool) string {
	n := field.Name
	t, found := mapGQLTypes[field.Type]
	if !found {
		t = "TODO # " + field.Type
	}

	if strings.HasSuffix(strings.ToLower(n), "id") && t == "String!" {
		t = "ID!"
	}

	if mustBeOptional {
		t = strings.TrimSuffix(t, "!")
	}

	return n + ": " + t
}

func isReadOnlyField(name string) bool {
	return (name == "ID" ||
		name == "CreatedAt" ||
		name == "UpdatedAt")
}

func GraphqlFields(
	fields []fieldDef,
	skipReadOnly bool,
) []fieldDef {
	var result []fieldDef

	for _, field := range fields {
		if skipReadOnly && isReadOnlyField(field.Name) {
			continue
		}

		result = append(
			result,
			fieldDef{
				Name: strcase.ToLowerCamel(strings.ReplaceAll(field.Name, "ID", "Id")),
				Type: field.Type,
			},
		)
	}

	return result
}
