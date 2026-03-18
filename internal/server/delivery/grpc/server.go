// Package grpc реализует gRPC-сервер для обработки запросов.
package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

// Service определяет интерфейс бизнес-логики.
type Service interface {
	Register(ctx context.Context, user *types.User) (token string, err error)
	Login(ctx context.Context, user *types.User) (token string, err error)

	GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error)
	CreateSecret(ctx context.Context, secretPayload *types.SecretPayload, userID int64) error
	GetSecret(ctx context.Context, secretID, userID int64) (*types.SecretPayload, error)
	UpdateSecret(ctx context.Context, updatedSecret *types.UpdatedSecret, userID int64) error
	DeleteSecret(ctx context.Context, secretID, userID int64) error
}

type options struct {
	gRPCAddress string
	jwtSecret   []byte
	tlsCertFile string
	tlsKeyFile  string
}

// Option определяет функцию настройки gRPC-сервера.
type Option func(*options)

// WithGRPCAddress устанавливает адрес gRPC-сервера.
func WithGRPCAddress(gRPCAddress string) Option {
	return func(o *options) {
		o.gRPCAddress = gRPCAddress
	}
}

// WithJWTSecret устанавливает секретный ключ для JWT.
func WithJWTSecret(jwtSecret []byte) Option {
	return func(o *options) {
		o.jwtSecret = jwtSecret
	}
}

// WithTLSCertFile устанавливает путь к TLS-сертификату.
func WithTLSCertFile(tlsCertFile string) Option {
	return func(o *options) {
		o.tlsCertFile = tlsCertFile
	}
}

// WithTLSKeyFile устанавливает путь к приватному ключу TLS.
func WithTLSKeyFile(tlsKeyFile string) Option {
	return func(o *options) {
		o.tlsKeyFile = tlsKeyFile
	}
}

// Server представляет gRPC-сервер.
type Server struct {
	pb.UnimplementedKeeperServer
	address string
	server  *grpc.Server
	service Service
	logger  zerolog.Logger
}

// New создаёт новый gRPC-сервер.
func New(service Service, logger zerolog.Logger, opts ...Option) (*Server, error) {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(LoggerInterceptor(logger), AuthInterceptor(o.jwtSecret)),
	}

	if o.tlsCertFile != "" && o.tlsKeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(o.tlsCertFile, o.tlsKeyFile)
		if err != nil {
			return nil, errors.Wrap(err, "ошибка загрузки TLS сертификатов")
		}
		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	return &Server{
		address: o.gRPCAddress,
		server:  grpc.NewServer(grpcOpts...),
		service: service,
		logger:  logger,
	}, nil
}

// Run запускает gRPC-сервер.
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return errors.Wrap(err, "не удалось открыть tcp-сокет для gRPC-сервера")
	}

	pb.RegisterKeeperServer(s.server, s)

	s.logger.Info().Msgf("gRPC-сервер запущен на %s", s.address)

	if err := s.server.Serve(listener); err != nil {
		return errors.Wrap(err, "gRPC-сервер завершился с ошибкой")
	}

	return nil
}

// Shutdown gracefully останавливает gRPC-сервер.
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		s.logger.Info().Msg("Запущен graceful shutdown gRPC-сервера")
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		s.logger.Info().Msg("gRPC-сервер корректно остановлен")
		return nil
	}
}
