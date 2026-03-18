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

// GetUserByLogin возвращает пользователя по логину.
func (p *Postgres) GetUserByLogin(ctx context.Context, login string) (user *types.User, err error) {
	const q = `select id, login, password_hash, encrypted_user_key
from users
where login = :login`

	args := map[string]any{
		"login": login,
	}

	user = &types.User{}
	if err := dbutils.NamedGet(ctx, p.db, user, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.ErrInvalidCredentials
		}
		return nil, err
	}

	return user, nil
}
