// Package postgres реализует хранилище данных на базе PostgreSQL.
package postgres

import (
	"context"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
)

// GetSecrets возвращает список всех секретов пользователя.
func (p *Postgres) GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error) {
	const q = `select id, name, type, encrypted_data, metadata, user_id, created_at, updated_at
from secrets
where user_id = :user_id
order by id`

	args := map[string]any{
		"user_id": userID,
	}
	return secrets, dbutils.NamedSelect(ctx, p.db, &secrets, q, args)
}
