package postgres

import (
	"context"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/dbutils"
)

func (p *Postgres) UpdateSecret(ctx context.Context, secret *types.Secret) error {
	const q = `update secrets
set name           = :name,
    type           = :type,
    encrypted_data = :encrypted_data,
    metadata       = :metadata,
    updated_at     = :updated_at
where id = :id`

	return dbutils.NamedExec(ctx, p.db, q, secret)
}
