package {{.PkgName}}_test

import (
	"testing"

	"{{.ModuleName}}/db"
	"{{.ModuleName}}/pkg/{{.PkgName}}"
	"{{.ModuleName}}/tests"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func Test{{.NamePlural}}(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "{{.NamePlural}} Suite")
}

var DB db.DB
var _ = tests.SetupSuite(&DB)
var subject *{{.PkgName}}.{{.Name}}Service

var _ = BeforeEach(func() {
	subject = {{.PkgName}}.NewService(DB)
})
