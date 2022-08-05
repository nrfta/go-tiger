package {{.PkgName}}_test

import (
	"context"

	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/tests"
	"{{.ModuleName}}/tests/factories"

	"github.com/google/uuid"
	"github.com/nrfta/go-platform-security-policy/pkg/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} Service Test", func() {
	Describe("#Update", func() {
		var ctx context.Context
		var input gql_types.{{.Name}}UpdateInput
		var record *models.{{.Name}}

		BeforeEach(func() {
			ctx = tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.IdentifiedDomains.SupportLevel3,
			)

      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// {{ToLowerCamel $field.Name}} := {{$field.Type}}{{end}}
      {{ end }}
      input = gql_types.{{.Name}}UpdateInput{
        // TODO
        {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// {{$field.Name}}: &{{ToLowerCamel $field.Name}},{{end}}
        {{ end }}
      }
    })

		BeforeEach(func() {
			record = factories.Create[*models.{{.Name}}](
				DB,
				factories.{{.Name}},

        // TODO
        {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// factory.Use({{$field.Type}}).For("{{$field.Name}}"),{{end}}
        {{ end }}
			)
		})

		PIt("fails to update due to missing permission", func() {
			ctx := tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.NamedDomains.User,
			)

			_, err := subject.Update(ctx, record.ID, input)
			Expect(err).To(HaveOccurred())
		})

		It("updates the record", func() {
			result, err := subject.Update(ctx, record.ID, input)
			Expect(err).To(Succeed())

      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}Expect(result.{{$field.Name}}).To(Equal(*input.{{$field.Name}})){{end}}
      {{ end }}
		})
	})
})
