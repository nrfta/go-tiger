package generators

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/nrfta/go-tiger/pkg/generators/blueprints"
)

func NewGenerator[T any](
	rootDest string,
	data T,
	funcMap template.FuncMap,
) *Generator[T] {
	return &Generator[T]{
		Root:    rootDest,
		Data:    data,
		FuncMap: funcMap,
	}
}

type Generator[T any] struct {
	Root    string
	Data    T
	FuncMap template.FuncMap

	createdFiles  []string
	modifiedFiles []string
}

func (g *Generator[T]) AppendTemplateToFile(templatePath string, filePathToAppend string) {
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
		g.FuncMap,
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

func (g *Generator[T]) RenderTemplateToFile(templatePath string, destFilePath string) {
	content, err := blueprints.F.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Unable to read template file (%s): %s", templatePath, err.Error())
	}

	rendered, err := renderTemplate(
		templatePath,
		string(content),
		g.Data,
		g.FuncMap,
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

func (g *Generator[T]) PrintSummary() {
	fmt.Println("Created:")
	for _, file := range g.createdFiles {
		fmt.Println("  " + file)
	}

	fmt.Println("Modified:")
	for _, file := range g.modifiedFiles {
		fmt.Println("  " + file)
	}
}
