// Package postgres реализует хранилище данных на базе PostgreSQL.
package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

// GetUserByID возвращает пользователя по ID.
func (p *Postgres) GetUserByID(ctx context.Context, id int64) (user *types.User, err error) {
	const q = `select id, login, password_hash, encrypted_user_key
from users
where id = :id`

	args := map[string]any{
		"id": id,
	}

	user = &types.User{}
	if err := dbutils.NamedGet(ctx, p.db, user, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}
