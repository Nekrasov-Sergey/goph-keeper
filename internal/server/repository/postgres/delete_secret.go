// Package postgres реализует хранилище данных на базе PostgreSQL.
package postgres

import (
	"context"

	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
)

// DeleteSecret удаляет секрет по ID.
func (p *Postgres) DeleteSecret(ctx context.Context, secretID, userID int64) error {
	const q = `delete from secrets where id = :id and user_id = :user_id`

	args := map[string]any{
		"id":      secretID,
		"user_id": userID,
	}

	return dbutils.NamedExec(ctx, p.db, q, args)
}
