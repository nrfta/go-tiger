package scaffolding

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
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
	g.CreateFactory()
	g.CreateResolver()
	g.AddGraphqlQueries()
	g.AddGraphqlMutations()
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

func (g *Generator) createPkgFolder() string {
	p := path.Join(g.Root, "pkg", g.Data.PkgName)

	if err := os.MkdirAll(p, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	return p
}

func (g *Generator) CreatePkg() {
	pkgPath := g.createPkgFolder()

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

			destFileName := strings.ReplaceAll(
				strings.ReplaceAll(file.Name(), ".tpl", ""),
				"pkgName",
				g.Data.PkgName,
			)

			err = os.WriteFile(
				path.Join(pkgPath, destFileName),
				[]byte(*rendered),
				os.ModePerm,
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func (g *Generator) AddGraphqlQueries() {
	g.appendTemplateToFile("misc/query.graphql.tpl", "pkg/schemas/query.graphql")
}

func (g *Generator) AddGraphqlMutations() {
	g.appendTemplateToFile("misc/mutation.graphql.tpl", "pkg/schemas/mutation.graphql")
}

func (g *Generator) appendTemplateToFile(templatePath string, filePathToAppend string) {
	var file, err = os.OpenFile(
		path.Join(g.Root, filePathToAppend),
		os.O_RDWR,
		os.ModePerm,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lastLineSize := 0
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}

		lastLineSize = len(line)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	offset := fileInfo.Size() - int64(lastLineSize+1)

	content, err := blueprints.F.ReadFile(templatePath)
	if err != nil {
		log.Fatal(err)
	}

	rendered, err := renderTemplate(
		"template",
		string(content),
		g.Data,
	)
	if err != nil {
		log.Fatal(err)
	}

	_, err = file.WriteAt(
		[]byte(*rendered),
		offset,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Generator) CreateFactory() {
	content, err := blueprints.F.ReadFile("misc/factory.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	rendered, err := renderTemplate(
		"factory.go",
		string(content),
		g.Data,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(
		path.Join(
			g.Root,
			"tests",
			"factories",
			strcase.ToSnake(g.Data.Name)+".go",
		),
		[]byte(*rendered),
		os.ModePerm,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Generator) CreateResolver() {
	content, err := blueprints.F.ReadFile("misc/resolver.go.tpl")
	if err != nil {
		log.Fatal(err)
	}

	rendered, err := renderTemplate(
		"resolver.go",
		string(content),
		g.Data,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(
		path.Join(
			g.Root,
			"pkg",
			"resolvers",
			strcase.ToSnake(g.Data.Name)+".go",
		),
		[]byte(*rendered),
		os.ModePerm,
	)
	if err != nil {
		log.Fatal(err)
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
	"ToLowerCamel":    strcase.ToLowerCamel,
	"GraphqlFields":   GraphqlFields,
	"GraphqlField":    GraphqlField,
	"isReadOnlyField": isReadOnlyField,
}

var mapGQLTypes = map[string]string{
	"string":            "String!",
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
