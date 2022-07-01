package {{.PkgName}}

import (
	"context"
	"database/sql"

	"{{.ModuleName}}/pkg/gql_types"
	"{{.ModuleName}}/pkg/models"

	"github.com/neighborly/go-errors"
	"github.com/neighborly/go-pghelpers"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (s *{{.Name}}Service) Update(
	ctx context.Context,
	id string,
	input gql_types.{{.Name}}UpdateInput,
) (*models.{{.Name}}, error) {
	record, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO
	// err = policy.{{.Name}}Update(ctx, id)
	// if err != nil {
	// 	return nil, errors.WithDisplayMessage(
	// 		err,
	// 		"Unable to update {{.Name}}, permission denied.",
	// 	)
	// }

  {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}
	if input.{{$field.Name}} != nil {
		record.{{$field.Name}} = *input.{{$field.Name}}
	}
  {{ end }}
  {{ end }}


	txErr := pghelpers.ExecInTx(s.db, func(tx *sql.Tx) bool {
		if _, err = record.Update(ctx, tx, boil.Infer()); err != nil {
			return false
		}

		return true
	})
	if txErr != nil {
		return nil, errors.Wrap(err, "failed to store db changes")
	}

	return record, err
}
