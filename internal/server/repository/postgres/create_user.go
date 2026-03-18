// Package postgres реализует хранилище данных на базе PostgreSQL.
package postgres

import (
	"context"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

// CreateUser создаёт нового пользователя в базе данных.
func (p *Postgres) CreateUser(ctx context.Context, user *types.User) (userID int64, err error) {
	const q = `insert into users (login, password_hash, encrypted_user_key)
values (:login, :password_hash, :encrypted_user_key)
returning id`

	if err := dbutils.NamedGet(ctx, p.db, &userID, q, user); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, errcodes.ErrLoginAlreadyExists
		}
		return 0, errors.Wrapf(err, "не удалось создать пользователя %s", user.Login)
	}

	return userID, nil
}
