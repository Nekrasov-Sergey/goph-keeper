package client

import (
	"github.com/rs/zerolog"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/config"
	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

type Client struct {
	config     *config.ClientConfig
	grpcClient pb.KeeperClient
	logger     zerolog.Logger
	token      string
}

func New(
	config *config.ClientConfig,
	grpcClient pb.KeeperClient,
) *Client {
	return &Client{
		config:     config,
		grpcClient: grpcClient,
	}
}
