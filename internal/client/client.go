// Package client реализует CLI-клиент для взаимодействия с сервером.
package client

import (
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

// Client представляет CLI-клиент для работы с секретами.
type Client struct {
	GRPCClient pb.KeeperClient
	token      string
}

// New создаёт новый экземпляр клиента.
func New(grpcClient pb.KeeperClient) *Client {
	return &Client{
		GRPCClient: grpcClient,
	}
}
