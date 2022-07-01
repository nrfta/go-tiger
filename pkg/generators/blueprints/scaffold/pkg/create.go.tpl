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

func (s *{{.Name}}Service) Create(
	ctx context.Context,
	input gql_types.{{.Name}}CreateInput,
) (*models.{{.Name}}, error) {
	record := &models.{{.Name}}{
    {{ range $index, $field := .Fields }}{{if not (isReadOnlyField $field.Name)}}{{$field.Name}}: input.{{$field.Name}},{{end}}
    {{ end }}
	}

	var err error
	txErr := pghelpers.ExecInTx(s.db, func(tx *sql.Tx) bool {
		if err = record.Insert(ctx, tx, boil.Infer()); err != nil {
			return false
		}

		// TODO
		// err = policy.{{.Name}}Create(ctx, record.ID)
		// if err != nil {
		// 	err = errors.WithDisplayMessage(
		// 		err,
		// 		"Unable to create {{.Name}}, permission denied.",
		// 	)

		// 	return false
		// }

		return true
	})
	if txErr != nil {
		return nil, errors.Wrap(err, "failed to store db record")
	}

	return record, err
}
