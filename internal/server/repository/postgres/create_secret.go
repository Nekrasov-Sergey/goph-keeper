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

func (p *Postgres) CreateSecret(ctx context.Context, secret *types.Secret) error {
	const q = `insert into secrets (name, type, encrypted_data, metadata, user_id, created_at)
values (:name, :type, :encrypted_data, :metadata, :user_id, :created_at)`

	if err := dbutils.NamedExec(ctx, p.db, q, secret); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errcodes.ErrSecretNameAlreadyExists
		}
		return errors.Wrapf(err, "не удалось создать секрет %s", secret.Name)
	}

	return nil
}
