package grpc

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

type KeeperClient struct {
	Client pb.KeeperClient
	conn   *grpc.ClientConn
}

func New(gRPCAddress string) (*KeeperClient, error) {
	conn, err := grpc.NewClient(
		gRPCAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, errors.Wrapf(err, "ошибка подключения к gRPC серверу по адресу %s", gRPCAddress)
	}

	return &KeeperClient{
		Client: pb.NewKeeperClient(conn),
		conn:   conn,
	}, nil
}

func (c *KeeperClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "ошибка закрытия gRPC соединения")
	}
	return nil
}
