package {{.PkgName}}

import (
	"context"

	"{{.ModuleName}}/db"
	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/pkg/utils"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/neighborly/go-errors"
	"github.com/nrfta/go-paging"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

var (
	// Err{{.Name}}NotFound used when a {{ToLowerCamel .Name}} is not found
	Err{{.Name}}NotFound = errors.WithDisplayMessage(
		errors.NotFound.New("unable to find {{ToLowerCamel .Name}}"),
		"{{.Name}} Not Found",
	)
)

// {{.Name}}Service implements the public api to interact with {{ToLowerCamel .NamePlural}}
type {{.Name}}Service struct {
	db db.DB
}

// {{.Name}}Connection is used for paginated result
type {{.Name}}Connection struct {
	Edges    []*{{.Name}}Edge
	PageInfo *paging.PageInfo
}

// {{.Name}}Edge is used for returning a cursor and {{.Name}}
type {{.Name}}Edge struct {
	Cursor *string
	Node   *models.{{.Name}}
}

// NewService creates an instance of the service containing the db
func NewService(
	db db.DB,
) *{{.Name}}Service {
	return &{{.Name}}Service{
		db: db,
	}
}

// Get a single {{ToLowerCamel .Name}} by ID
func (s *{{.Name}}Service) Get(
	ctx context.Context,
	id string,
) (*models.{{.Name}}, error) {
	{{ToLowerCamel .Name}}, err := models.Find{{.Name}}(ctx, s.db, id)

	if err != nil {
		return nil, Err{{.Name}}NotFound
	}

	return {{ToLowerCamel .Name}}, nil
}

// GetAll returns all {{ToLowerCamel .NamePlural}}
func (s *{{.Name}}Service) GetAll(
	ctx context.Context,
	mods ...qm.QueryMod,
) ([]*models.{{.Name}}, error) {
	{{ToLowerCamel .NamePlural}}, err := models.{{.NamePlural}}(mods...).All(ctx, s.db)
	if err != nil {
		return nil, err
	}

	return {{ToLowerCamel .NamePlural}}, err
}

// GetAllPaginated returns all {{ToLowerCamel .NamePlural}} with pagination
func (s *{{.Name}}Service) GetAllPaginated(
	ctx context.Context,
	page *paging.PageArgs,
	mods ...qm.QueryMod,
) (*{{.Name}}Connection, error) {
	totalCount, err := models.{{.NamePlural}}(mods...).Count(ctx, s.db)
	if err != nil {
		return &{{.Name}}Connection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	paginator := paging.NewOffsetPaginator(page, totalCount)
	mods = append(mods, paginator.QueryMods()...)

	records, err := models.{{.NamePlural}}(mods...).All(ctx, s.db)
	if err != nil {
		return &{{.Name}}Connection{
			PageInfo: paging.NewEmptyPageInfo(),
		}, err
	}

	result := &{{.Name}}Connection{
		PageInfo: &paginator.PageInfo,
	}

	for i, row := range records {
		result.Edges = append(result.Edges, &{{.Name}}Edge{
			Cursor: paging.EncodeOffsetCursor(paginator.Offset + i + 1),
			Node:   row,
		})
	}
	return result, nil
}

func (s *{{.Name}}Service) QueryModsForFilter(
	ctx context.Context,
	filter *gql_types.{{.Name}}Filter,
) []qm.QueryMod {
	var mods []qm.QueryMod

	if filter == nil {
		return mods
	}

	if filter.Ids != nil && len(filter.Ids) > 0 {
		mods = append(mods, models.{{.Name}}Where.ID.IN(filter.Ids))
	}

	return mods
}

func (s *{{.Name}}Service) NewLoaderByID() *dataloader.Loader[string, *models.{{.Name}}] {
	cache := dataloader.NewCache[string, *models.{{.Name}}]()
	loader := dataloader.NewBatchedLoader(
		utils.BatchedLoaderFn(
			func(ctx context.Context, IDs []string) ([]*models.{{.Name}}, error) {
				return s.GetAll(ctx, models.{{.Name}}Where.ID.IN(IDs))
			},
			func(obj *models.{{.Name}}) string {
				return obj.ID
			},
			Err{{.Name}}NotFound,
		),
		dataloader.WithCache[string, *models.{{.Name}}](cache),
	)

	return loader
}
