package {{.PkgName}}_test

import (
	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/tests"

	"github.com/google/uuid"
	"github.com/nrfta/go-platform-security-policy/pkg/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} Service Test", func() {
	Describe("#Create", func() {
		input := gql_types.{{.Name}}CreateInput{
      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// {{$field.Name}}: {{$field.Type}},{{end}}
      {{ end }}
		}

		PIt("fails to create due to missing permission", func() {
			ctx := tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.NamedDomains.User,
			)

			_, err := subject.Create(ctx, input)
			Expect(err).To(HaveOccurred())
		})

		It("creates the record", func() {
			ctx := tests.ContextWithCurrentUserInfo(
				uuid.NewString(),
				nil,
				policy.IdentifiedDomains.SupportLevel3,
			)

			result, err := subject.Create(ctx, input)
			Expect(err).To(Succeed())

      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}Expect(result.{{$field.Name}}).To(Equal(input.{{$field.Name}})){{end}}
      {{ end }}
		})
	})
})
