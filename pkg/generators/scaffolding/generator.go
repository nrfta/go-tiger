package scaffolding

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/nrfta/go-tiger/pkg/generators"
	"github.com/nrfta/go-tiger/pkg/generators/blueprints"
	"github.com/volatiletech/inflect"
	"golang.org/x/mod/modfile"
)

type Scaffold struct {
	Gen *generators.Generator[*Data]
}

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

	s := &Scaffold{
		Gen: generators.NewGenerator(
			rootPath,
			data,
			funcMap,
		),
	}

	s.CreatePkg()
	s.CreateFactory()
	s.CreateResolver()
	s.AddGraphqlQueries()
	s.AddGraphqlMutations()

	s.Gen.PrintSummary()

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

func (s *Scaffold) CreatePkg() {
	pkgPath := s.createPkgFolder()

	files, err := blueprints.F.ReadDir("scaffold/pkg")
	if err != nil {
		log.Fatalf("Unable to read template dir (scaffold/pkg): %s", err.Error())
	}

	for _, file := range files {
		if !file.IsDir() {
			destFileName := strings.ReplaceAll(
				strings.ReplaceAll(file.Name(), ".tpl", ""),
				"pkgName",
				s.Gen.Data.PkgName,
			)

			s.Gen.RenderTemplateToFile("scaffold/pkg/"+file.Name(), pkgPath+"/"+destFileName)
		}
	}
}

func (g *Scaffold) CreateFactory() {
	g.Gen.RenderTemplateToFile(
		"scaffold/misc/factory.go.tpl",
		"tests/factories/"+strcase.ToSnake(g.Gen.Data.Name)+".go",
	)
}

func (g *Scaffold) CreateResolver() {
	g.Gen.RenderTemplateToFile(
		"scaffold/misc/resolver.go.tpl",
		"pkg/resolvers/"+strcase.ToSnake(g.Gen.Data.Name)+".go",
	)
}

func (g *Scaffold) AddGraphqlQueries() {
	g.Gen.AppendTemplateToFile("scaffold/misc/query.graphql.tpl", "pkg/schemas/query.graphql")
}

func (g *Scaffold) AddGraphqlMutations() {
	g.Gen.AppendTemplateToFile("scaffold/misc/mutation.graphql.tpl", "pkg/schemas/mutation.graphql")
}

func (g *Scaffold) createPkgFolder() string {
	p := "pkg/" + g.Gen.Data.PkgName

	if err := os.MkdirAll(path.Join(g.Gen.Root, p), 0644); err != nil {
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
