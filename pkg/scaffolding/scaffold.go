package scaffolding

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/nrfta/go-tiger/pkg/scaffolding/blueprints"
	"github.com/volatiletech/inflect"
	"golang.org/x/mod/modfile"
)

func Process(rootPath, filePath string) {
	m := ParseModel(filePath)

	namePlural := inflect.Pluralize(m.EntityName)
	data := &Data{
		ModuleName: getModuleName(rootPath),
		Name:       m.EntityName,
		NamePlural: namePlural,
		PkgName:    strcase.ToLowerCamel(namePlural),
		Fields:     m.Fields,
	}

	g := &Generator{
		Root: rootPath,
		Data: data,
	}

	g.CreatePkg()
}

type Data struct {
	ModuleName string
	PkgName    string
	Name       string
	NamePlural string
	Fields     []fieldDef
}

type Generator struct {
	Root string
	Data *Data
}

func (g *Generator) createPkgFolder() {
	p := path.Join(g.Root, "pkg", g.Data.PkgName)

	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

func (g *Generator) CreatePkg() {
	g.createPkgFolder()

	files, err := blueprints.F.ReadDir("pkg")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			file.Name()

			content, err := blueprints.F.ReadFile("pkg/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			rendered, err := renderTemplate(
				file.Name(),
				string(content),
				g.Data,
			)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(*rendered)

		}
	}

}

func getModuleName(rootPath string) string {
	goModBytes, err := ioutil.ReadFile(path.Join(rootPath, "go.mod"))
	if err != nil {
		log.Fatal(err)
	}

	return modfile.ModulePath(goModBytes)
}

var funcMap template.FuncMap = template.FuncMap{
	"ToSnake":      strcase.ToSnake,
	"ToCamel":      strcase.ToCamel,
	"ToLowerCamel": strcase.ToLowerCamel,
}

func renderTemplate(
	name string,
	input string,
	data interface{},
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
