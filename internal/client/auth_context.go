// Package client реализует CLI-клиент для взаимодействия с сервером.
package client

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// AuthContext добавляет JWT-токен в контекст запроса.
func (c *Client) AuthContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + c.token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
