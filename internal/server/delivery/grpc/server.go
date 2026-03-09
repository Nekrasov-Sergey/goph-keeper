package grpc

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"

	pb "github.com/Nekrasov-Sergey/goph-keeper/internal/proto"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/types"
)

type Service interface {
	Register(ctx context.Context, user *types.User) (token string, err error)
	Login(ctx context.Context, user *types.User) (token string, err error)

	GetSecrets(ctx context.Context, userID int64) (secrets []types.Secret, err error)
	CreateSecret(ctx context.Context, secretInput *types.SecretInput, userID int64) error
	GetSecret(ctx context.Context, secretID, userID int64) (*types.SecretInput, error)
	UpdateSecret(ctx context.Context, updatedSecret *types.UpdatedSecret, userID int64) error
	DeleteSecret(ctx context.Context, secretID, userID int64) error
}

type options struct {
	gRPCAddress string
	jwtSecret   []byte
}

type Option func(*options)

func WithGRPCAddress(gRPCAddress string) Option {
	return func(o *options) {
		o.gRPCAddress = gRPCAddress
	}
}

func WithJWTSecret(jwtSecret []byte) Option {
	return func(o *options) {
		o.jwtSecret = jwtSecret
	}
}

type Server struct {
	pb.UnimplementedKeeperServer
	address string
	server  *grpc.Server
	service Service
	logger  zerolog.Logger
}

func New(service Service, logger zerolog.Logger, opts ...Option) (*Server, error) {
	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	return &Server{
		address: o.gRPCAddress,
		server:  grpc.NewServer(grpc.ChainUnaryInterceptor(LoggerInterceptor(logger), AuthInterceptor(o.jwtSecret))),
		service: service,
		logger:  logger,
	}, nil
}

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
