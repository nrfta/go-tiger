package resolvers

import (
	"context"

	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"
	"{{.ModuleName}}/pkg/{{ToLowerCamel .Name}}s"

	"github.com/nrfta/go-paging"
)

func (r *queryResolver) {{.Name}}ByID(
	ctx context.Context,
	id string,
) (*models.{{.Name}}, error) {
	return r.Services.{{.Name}}.Get(ctx, id)
}

func (r *queryResolver) {{.NamePlural}}(
	ctx context.Context,
	page *paging.PageArgs,
	filter *gql_types.{{.Name}}Filter,
) (*{{ToLowerCamel .Name}}s.{{.Name}}Connection, error) {
	return r.Services.{{.Name}}.GetAllPaginated(
		ctx,
		page,
		r.Services.{{.Name}}.QueryModsForFilter(ctx, filter)...,
	)
}

func (r *mutationResolver) Create{{.Name}}(
	ctx context.Context,
	input gql_types.{{.Name}}CreateInput,
) (*models.{{.Name}}, error) {
	return r.Services.{{.Name}}.Create(ctx, input)
}

func (r *mutationResolver) Update{{.Name}}(
	ctx context.Context,
	id string,
	input gql_types.{{.Name}}UpdateInput,
) (*models.{{.Name}}, error) {
	return r.Services.{{.Name}}.Update(ctx, id, input)
}

func (r *mutationResolver) Delete{{.Name}}(
	ctx context.Context,
	id string,
) (*models.{{.Name}}, error) {
	return r.Services.{{.Name}}.Delete(ctx, id)
}

// func (r *rootResolver) {{.Name}}() {{.Name}}Resolver {
// 	return &{{ToLowerCamel .Name}}Resolver{Services: r.Services}
// }
// 
// type {{ToLowerCamel .Name}}Resolver struct {
// 	Services *ServiceRegistry
// }
