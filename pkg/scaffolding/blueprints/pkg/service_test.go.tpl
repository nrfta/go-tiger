package {{.PkgName}}_test

import (
	"context"

	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/tests/factories"

	"github.com/google/uuid"
	"github.com/kolach/go-factory"
	"github.com/nrfta/go-paging"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("{{.NamePlural}} Service Test", func() {
	Describe("#Get", func() {
		It("returns error when record is not found", func() {
			randomID := uuid.NewString()
			_, err := subject.Get(context.Background(), randomID)
			Expect(err).To(HaveOccurred())
		})

		It("returns the record", func() {
			record := factories.Create(
				DB,
				factories.{{.Name}},
			).(*models.{{.Name}})

			result, err := subject.Get(context.Background(), record.ID)
			Expect(err).To(Succeed())
			Expect(result.ID).To(Equal(record.ID))
		})
	})

	Describe("Pagination, Filters and Dataloader", func() {
		var (
			ctx     context.Context
			record1 *models.{{.Name}}
			record2 *models.{{.Name}}
			record3 *models.{{.Name}}
		)

		BeforeEach(func() {
			ctx = context.Background()

			record1 = factories.Create(
				DB,
				factories.{{.Name}},

				// factory.Use(false).For("SomeField"),
			).(*models.{{.Name}})

			record2 = factories.Create(
				DB,
				factories.{{.Name}},
			).(*models.{{.Name}})

			record3 = factories.Create(
				DB,
				factories.{{.Name}},
			).(*models.{{.Name}})
		})

		Describe("#GetAllPaginated", func() {
			It("returns no error when no records are found", func() {
				mods := subject.QueryModsForFilter(
					context.Background(),
					&gql_types.{{.Name}}Filter{Ids: []string{uuid.NewString()}},
				)
				records, err := subject.GetAllPaginated(context.Background(), nil, mods...)
				Expect(err).To(Succeed())
				Expect(len(records.Edges)).To(Equal(0))
			})

			It("should return all records, paginated", func() {
				limit := 2
				result, err := subject.GetAllPaginated(
					context.Background(),
					&paging.PageArgs{First: &limit},
				)
				Expect(err).To(Succeed())
				Expect(len(result.Edges)).To(Equal(2))

				count, _ := result.PageInfo.TotalCount()
				Expect(*count).To(Equal(3))
				haveNext, _ := result.PageInfo.HasNextPage()
				Expect(haveNext).To(BeTrue())

				page2, err := subject.GetAllPaginated(
					context.Background(),
					&paging.PageArgs{
						First: &limit,
						After: result.Edges[1].Cursor,
					},
				)
				Expect(err).To(Succeed())
				Expect(page2.Edges).To(HaveLen(1))
				haveNext, _ = page2.PageInfo.HasNextPage()
				Expect(haveNext).To(BeFalse())
			})
		})

		Describe("#QueryModsForFilter / #GetAll", func() {
			// ids: [ID!]
			It("should filter by IDs", func() {
				ids := []string{record1.ID, record3.ID}
				mods := subject.QueryModsForFilter(
					ctx,
					&gql_types.{{.Name}}Filter{Ids: ids},
				)
				records, err := subject.GetAll(ctx, mods...)
				Expect(err).To(Succeed())
				Expect(records).To(HaveLen(2))
				for _, v := range records {
					Expect(v.ID).To(BeElementOf(ids))
				}

				ids = []string{record2.ID}
				mods = subject.QueryModsForFilter(
					ctx,
					&gql_types.{{.Name}}Filter{Ids: ids},
				)
				records, err = subject.GetAll(ctx, mods...)
				Expect(err).To(Succeed())
				Expect(records).To(HaveLen(1))
				Expect(records[0].ID).To(Equal(record2.ID))
			})
		})

		Describe("#NewLoaderByID", func() {
			It("returns ordered records by id", func() {
				ids := []string{record3.ID, record1.ID}

				loader := subject.NewLoaderByID()

				records, errs := loader.LoadMany(ctx, ids)()
				Expect(errs).To(HaveLen(0))

				Expect(records).To(HaveLen(2))
				Expect(records[0].ID).To(Equal(ids[0]))
				Expect(records[1].ID).To(Equal(ids[1]))
			})
		})
	})
})
