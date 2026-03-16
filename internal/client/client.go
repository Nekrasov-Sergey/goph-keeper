package client

import (
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

type Client struct {
	GRPCClient pb.KeeperClient
	token      string
}

func New(grpcClient pb.KeeperClient) *Client {
	return &Client{
		GRPCClient: grpcClient,
	}
}
