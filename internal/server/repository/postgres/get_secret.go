package postgres

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/errcodes"
)

func (p *Postgres) GetSecret(ctx context.Context, secretID, userID int64) (secret *types.Secret, err error) {
	const q = `select id, name, type, encrypted_data, metadata, user_id, created_at, updated_at
from secrets
where id = :id and user_id = :user_id
order by created_at`

	args := map[string]any{
		"id":      secretID,
		"user_id": userID,
	}

	secret = &types.Secret{}
	if err := dbutils.NamedGet(ctx, p.db, secret, q, args); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errcodes.ErrSecretNotFound
		}
		return nil, err
	}

	return secret, nil
}
