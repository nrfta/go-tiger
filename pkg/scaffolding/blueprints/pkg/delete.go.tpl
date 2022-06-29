package {{.PkgName}}

import (
	"context"
	"database/sql"

	"{{.ModuleName}}/pkg/models"

	"github.com/neighborly/go-errors"
	"github.com/neighborly/go-pghelpers"
)

func (s *{{.Name}}Service) Delete(
	ctx context.Context,
	id string,
) (*models.{{.Name}}, error) {
	record, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO
	// err = policy.{{.Name}}Delete(ctx, id)
	// if err != nil {
	// 	return nil, errors.WithDisplayMessage(
	// 		err,
	// 		"Unable to delete {{.Name}}, permission denied.",
	// 	)
	// }

	txErr := pghelpers.ExecInTx(s.db, func(tx *sql.Tx) bool {
		if _, err = record.Delete(ctx, tx); err != nil {
			return false
		}

		return true
	})
	if txErr != nil {
		return nil, errors.Wrap(err, "failed to store db changes")
	}

	return record, err
}
