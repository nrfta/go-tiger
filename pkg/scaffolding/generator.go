package scaffolding

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

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

	fmt.Println("Created:")
	for _, file := range g.createdFiles {
		fmt.Println("  " + file)
	}

	fmt.Println("Modified:")
	for _, file := range g.modifiedFiles {
		fmt.Println("  " + file)
	}

	fmt.Println()
	fmt.Println()
	fmt.Println("What's Next?")
	fmt.Println("  Manually Register the Service in pkg/resolvers/service_registry.go")
	fmt.Println("  Add new pkg to autobind in gqlgen.yml")
	fmt.Println("  Fix TODOs")
	fmt.Println("  Implement Access Management")
	fmt.Println("  Review generated GraphQL:")
	fmt.Println("    - Remove any undesired field")
	fmt.Println("    - Use graphql type for foreign key fields if any")

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

	createdFiles  []string
	modifiedFiles []string
}

func (g *Generator) CreatePkg() {
	pkgPath := g.createPkgFolder()

	files, err := blueprints.F.ReadDir("pkg")
	if err != nil {
		log.Fatalf("Unable to read template dir (pkg): %s", err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			destFileName := strings.ReplaceAll(
				strings.ReplaceAll(file.Name(), ".tpl", ""),
				"pkgName",
				g.Data.PkgName,
			)

			g.renderTemplateToFile("pkg/"+file.Name(), pkgPath+"/"+destFileName)
		}
	}
}

func (g *Generator) CreateFactory() {
	g.renderTemplateToFile(
		"misc/factory.go.tpl",
		"tests/factories/"+strcase.ToSnake(g.Data.Name)+".go",
	)
}

func (g *Generator) CreateResolver() {
	g.renderTemplateToFile(
		"misc/resolver.go.tpl",
		"pkg/resolvers/"+strcase.ToSnake(g.Data.Name)+".go",
	)
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
		0644,
	)
	if err != nil {
		log.Fatalf("Unable to read file to append (%s): %s", filePathToAppend, err.Error())
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
		log.Fatalf("Unable to get file stat (%s): %s", filePathToAppend, err.Error())
	}
	offset := fileInfo.Size() - int64(lastLineSize+1)

	content, err := blueprints.F.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Unable to read template (%s): %s", templatePath, err.Error())
	}

	rendered, err := renderTemplate(
		"template",
		string(content),
		g.Data,
	)
	if err != nil {
		log.Fatalf("Unable to render template (%s): %s", templatePath, err.Error())
	}

	_, err = file.WriteAt(
		[]byte(*rendered),
		offset,
	)
	if err != nil {
		log.Fatalf("Unable to write at the end of the file (%s): %s", filePathToAppend, err.Error())
	}

	// Save file changes.
	err = file.Sync()
	if err != nil {
		log.Fatalf("Unable to save the file (%s): %s", filePathToAppend, err.Error())
	}

	g.modifiedFiles = append(g.modifiedFiles, filePathToAppend)
}

func (g *Generator) renderTemplateToFile(templatePath string, destFilePath string) {
	content, err := blueprints.F.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Unable to read template file (%s): %s", templatePath, err.Error())
	}

	rendered, err := renderTemplate(
		templatePath,
		string(content),
		g.Data,
	)
	if err != nil {
		log.Fatalf("Unable to render template (%s): %s", templatePath, err.Error())
	}

	err = os.WriteFile(
		path.Join(g.Root, destFilePath),
		[]byte(*rendered),
		0644,
	)
	if err != nil {
		log.Fatalf("Unable to write file (%s): %s", destFilePath, err.Error())
	}

	g.createdFiles = append(g.createdFiles, destFilePath)
}

func (g *Generator) createPkgFolder() string {
	p := "pkg/" + g.Data.PkgName

	if err := os.MkdirAll(path.Join(g.Root, p), 0644); err != nil {
		log.Fatalf("Unable to folder (%s): %s", p, err.Error())
	}

	return p
}

func getModuleName(rootPath string) string {
	goModBytes, err := ioutil.ReadFile(path.Join(rootPath, "go.mod"))
	if err != nil {
		log.Fatalf("Unable to find module name: %s", err.Error())
	}

	return modfile.ModulePath(goModBytes)
}
