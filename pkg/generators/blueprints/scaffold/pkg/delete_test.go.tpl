package {{.PkgName}}_test

import (
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/tests"
	"{{.ModuleName}}/tests/factories"

	"github.com/google/uuid"
	"github.com/nrfta/go-platform-security-policy/pkg/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} Service Test", func() {
	Describe("#Delete", func() {

		var record *models.{{.Name}}

		BeforeEach(func() {
			record = factories.Create(
				DB,
				factories.{{.Name}},

       // TODO
        {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// factory.Use({{$field.Type}}).For("{{$field.Name}}"),{{end}}
        {{ end }}
			).(*models.{{.Name}})
		})

		PIt("fails to delete due to missing permission", func() {
			ctx := tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.NamedDomains.User,
			)

			_, err := subject.Delete(ctx, record.ID)
			Expect(err).To(HaveOccurred())
		})

		It("deletes the record", func() {
			ctx := tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.IdentifiedDomains.SupportLevel3,
			)

			_, err := subject.Delete(ctx, record.ID)
			Expect(err).To(Succeed())
			Expect(record.Reload(ctx, DB)).To(HaveOccurred())
		})
	})
})
