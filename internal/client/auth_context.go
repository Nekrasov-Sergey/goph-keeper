package client

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func (c *Client) AuthContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"authorization": "Bearer " + c.token,
	})

	return metadata.NewOutgoingContext(ctx, md)
}
