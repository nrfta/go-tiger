package integration_test

import (
	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/pkg/{{ToLowerCamel .NamePlural}}"
	"{{.ModuleName}}/tests"
	"{{.ModuleName}}/tests/factories"

	"github.com/99designs/gqlgen/client"
	"github.com/google/uuid"
	"github.com/nrfta/go-platform-security-policy/pkg/policy"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} Integration Test", func() {
	// var ctx context.Context
	var token string
	var unauthorizedToken string

	supportLevel3UserID := "cc7157a7-1a78-40f8-81c8-645f106a91d1"

	BeforeEach(func() {
		unauthorizedToken, _ = tests.CreateUserToken(
			uuid.NewString(),
			[]string{},
			policy.NamedDomains.User,
		)

		token, _ = tests.CreateUserToken(
			supportLevel3UserID,
			[]string{},
			policy.IdentifiedDomains.SupportLevel3,
		)
	})

	Context("#{{ToLowerCamel .NamePlural}}", func() {
		var (
			gClient *client.Client
		)

		query := `
			query ($filter: {{.Name}}Filter) {
				{{ToLowerCamel .NamePlural}}(filter: $filter) {
					edges {
						node {
							id
						}
					}
				}
			}`

		var resp struct {
			{{.NamePlural}} {{ToLowerCamel .NamePlural}}.{{.Name}}Connection
		}

		BeforeEach(func() {
			gClient = newGraphqlClient(token)
		})

		It("returns the {{ToLowerCamel .NamePlural}}", func() {
			factories.CreateCount[*models.{{.Name}}](s.DB, factories.{{.Name}}, 2)

			gClient.MustPost(query, &resp)

			Expect(resp.{{.NamePlural}}.Edges).To(HaveLen(2))
		})

		//	ids: [ID!]
		It("filters by IDs", func() {
			records := factories.CreateCount[*models.{{.Name}}](s.DB, factories.{{.Name}}, 5)

			ids := []string{records[1].ID, records[3].ID}
			filter := gql_types.{{.Name}}Filter{Ids: ids}

			gClient.MustPost(
				query,
				&resp,
				client.Var("filter", filter),
			)

			Expect(resp.{{.NamePlural}}.Edges).To(HaveLen(2))
			for _, v := range resp.{{.NamePlural}}.Edges {
				Expect(v.Node.ID).To(BeElementOf(ids))
			}
		})

		It("fails when user is unauthenticated", func() {
			c := newGraphqlClient()

			err := c.Post(
				query,
				&resp,
				client.Var("filter", nil),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"no user info found","path":["{{ToLowerCamel .NamePlural}}"],"extensions":{"code":"UNAUTHENTICATED"}}]`))
		})

		It("fails when user is unauthorized", func() {
			factories.CreateCount[*models.{{.Name}}](s.DB, factories.{{.Name}}, 1)

			c := newGraphqlClient(unauthorizedToken)

			err := c.Post(
				query,
				&resp,
				client.Var("filter", nil),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"current user ID is not authorized to access the resource","path":["{{ToLowerCamel .NamePlural}}","edges",0,"node"],"extensions":{"code":"PERMISSION_DENIED"}}]`))
		})
	})

	Context("#{{ToLowerCamel .Name}}ByID", func() {
		query := `
			query($id: ID!) {
				{{ToLowerCamel .Name}}ById(id: $id) {
					id
				}
			}`

		var resp struct {
			{{.Name}}ByID struct {
				ID   string
			}
		}

		It("returns the {{ToLowerCamel .Name}}", func() {
			c := newGraphqlClient(token)

			{{ToLowerCamel .Name}} := factories.Create[*models.{{.Name}}](
				s.DB,
				factories.{{.Name}},
			)

			c.MustPost(
				query,
				&resp,
				client.Var("id", {{ToLowerCamel .Name}}.ID),
			)

			Expect(resp.{{.Name}}ByID.ID).To(Equal({{ToLowerCamel .Name}}.ID))
		})

		It("fails when unauthenticated", func() {
			c := newGraphqlClient()

			{{ToLowerCamel .Name}} := factories.Create[*models.{{.Name}}](
				s.DB,
				factories.{{.Name}},
			)

			err := c.Post(
				query,
				&resp,
				client.Var("id", {{ToLowerCamel .Name}}.ID),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"no user info found","path":["{{ToLowerCamel .Name}}ById"],"extensions":{"code":"UNAUTHENTICATED"}}]`))
		})

		It("fails when user is unauthorized", func() {
			c := newGraphqlClient(unauthorizedToken)

			{{ToLowerCamel .Name}} := factories.Create[*models.{{.Name}}](
				s.DB,
				factories.{{.Name}},
			)

			err := c.Post(
				query,
				&resp,
				client.Var("id", {{ToLowerCamel .Name}}.ID),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"current user ID is not authorized to access the resource","path":["{{ToLowerCamel .Name}}ById"],"extensions":{"code":"PERMISSION_DENIED"}}]`))
		})
	})

	Describe("#create{{.Name}}", func() {
		input := gql_types.{{.Name}}CreateInput{
      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// {{$field.Name}}: {{$field.Type}},{{end}}
      {{ end }}
		}

		mutation := `
			mutation($input: {{.Name}}CreateInput!) {
				create{{.Name}}(input: $input) {
					id
				}
			}`

		var resp struct {
			Create{{.Name}} struct {
				ID   string
			}
		}

		It("fails to create when unauthorized", func() {
			c := newGraphqlClient(unauthorizedToken)

			err := c.Post(
				mutation,
				&resp,
				client.Var("input", input),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"Unable to create {{.Name}}, permission denied.","path":["create{{.Name}}"],"extensions":{"code":"PERMISSION_DENIED"}}]`))
		})

		It("creates the {{ToLowerCamel .Name}}", func() {
			c := newGraphqlClient(token)

			err := c.Post(mutation, &resp,
				client.Var("input", input),
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Create{{.Name}}.ID).ToNot(BeEmpty())
		})
	})

	Describe("#update{{.Name}}", func() {
		var {{ToLowerCamel .Name}} *models.{{.Name}}

		BeforeEach(func() {
			{{ToLowerCamel .Name}} = factories.Create[*models.{{.Name}}](
				s.DB,
				factories.{{.Name}},
			)
		})

		name := "Updated Name"

		input := gql_types.{{.Name}}UpdateInput{
      // TODO
      {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}// {{$field.Name}}: &{{ToLowerCamel $field.Name}},{{end}}
      {{ end }}
		}

		mutation := `
			mutation($id: ID!, $input: {{.Name}}UpdateInput!) {
				update{{.Name}}(id: $id, input: $input) {
					id
				}
			}`

		var resp struct {
			Update{{.Name}} struct {
				ID   string
			}
		}

		It("fails to update when unauthorized", func() {
			c := newGraphqlClient(unauthorizedToken)

			err := c.Post(
				mutation,
				&resp,
				client.Var("id", {{ToLowerCamel .Name}}.ID),
				client.Var("input", input),
			)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(`[{"message":"Unable to update {{.Name}}, permission denied.","path":["update{{.Name}}"],"extensions":{"code":"PERMISSION_DENIED"}}]`))
		})

		It("updates the {{ToLowerCamel .Name}}", func() {
			c := newGraphqlClient(token)

			err := c.Post(
				mutation,
				&resp,
				client.Var("id", {{ToLowerCamel .Name}}.ID),
				client.Var("input", input),
			)

			Expect(err).ToNot(HaveOccurred())
			Expect(resp.Update{{.Name}}.ID).To(Equal({{ToLowerCamel .Name}}.ID))
		})
	})
})
