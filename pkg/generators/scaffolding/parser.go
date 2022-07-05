package scaffolding

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/nrfta/go-log"
)

type fieldDef struct {
	Name string
	Type string
}

func ParseModel(filePath string) *visitor {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		panic(err)
	}

	v := &visitor{}
	ast.Walk(v, file)

	return v
}

type visitor struct {
	EntityName string
	Fields     []fieldDef
}

func (v *visitor) Visit(node ast.Node) (w ast.Visitor) {
	switch t := node.(type) {
	case *ast.TypeSpec:
		if v.EntityName == "" {
			v.EntityName = t.Name.Name
			tStruct := t.Type.(*ast.StructType)

			for _, f := range tStruct.Fields.List {
				if len(f.Names) == 0 {
					continue
				}

				fieldName := f.Names[0].Name
				if fieldName == "R" {
					continue
				}
				if fieldName == "L" {
					continue
				}
				switch fieldType := f.Type.(type) {
				case *ast.Ident:
					v.Fields = append(
						v.Fields,
						fieldDef{
							Name: fieldName,
							Type: fieldType.Name,
						},
					)
				case *ast.SelectorExpr:
					packageName, ok := fieldType.X.(*ast.Ident)
					if !ok {
						log.Warnf("Unexpected field type : %v\n", f.Type)

						continue
					}
					selector := fieldType.Sel
					v.Fields = append(
						v.Fields,
						fieldDef{
							Name: fieldName,
							Type: packageName.Name + "." + selector.Name,
						},
					)

				default:
					log.Warnf("Unexpected field type : %v\n", f.Type)
					continue
				}
			}
		}
	}

	return v
}
