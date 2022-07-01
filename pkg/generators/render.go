package generators

import (
	"bytes"
	"text/template"
)

func renderTemplate(
	name string,
	input string,
	data interface{},
	funcMap template.FuncMap,
) (*string, error) {
	t, err := template.New(name).Funcs(funcMap).Parse(input)
	if err != nil {
		return nil, err
	}

	var content bytes.Buffer
	err = t.Execute(&content, data)
	if err != nil {
		return nil, err
	}

	contStr := content.String()
	return &contStr, nil
}
