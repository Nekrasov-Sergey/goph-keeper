// Package grpc реализует gRPC-клиент для взаимодействия с сервером.
package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
)

// options содержит параметры конфигурации gRPC-клиента.
type options struct {
	gRPCAddress string
	tlsCertFile string
}

// Option определяет функцию настройки gRPC-клиента.
type Option func(*options)

// WithGRPCAddress устанавливает адрес gRPC-сервера.
func WithGRPCAddress(gRPCAddress string) Option {
	return func(o *options) {
		o.gRPCAddress = gRPCAddress
	}
}

// WithTLSCertFile устанавливает путь к TLS-сертификату.
func WithTLSCertFile(tlsCertFile string) Option {
	return func(o *options) {
		o.tlsCertFile = tlsCertFile
	}
}

// KeeperClient представляет gRPC-клиент для работы с Keeper.
type KeeperClient struct {
	Client pb.KeeperClient
	conn   *grpc.ClientConn
}

// New создаёт новый gRPC-клиент с TLS.
func New(opts ...Option) (*KeeperClient, error) {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	cert, err := os.ReadFile(o.tlsCertFile)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка чтения TLS сертификата")
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(cert) {
		return nil, errors.New("ошибка добавления сертификата в pool")
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.NewClient(o.gRPCAddress, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrapf(err, "ошибка подключения к gRPC серверу по адресу %s", o.gRPCAddress)
	}

	return &KeeperClient{
		Client: pb.NewKeeperClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает gRPC-соединение.
func (c *KeeperClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return errors.Wrap(err, "ошибка закрытия gRPC соединения")
	}
	return nil
}
