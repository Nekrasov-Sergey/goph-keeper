package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/config"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/server/delivery/grpc"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/server/repository/postgres"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/server/service"
	buildinfo "github.com/Nekrasov-Sergey/goph-keeper/pkg/build_info"
	"github.com/Nekrasov-Sergey/goph-keeper/pkg/logger"
)

func main() {
	buildinfo.Print()
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Сервер завершился с ошибкой")
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	l := logger.New()

	cfg, err := config.NewServerConfig(l)
	if err != nil {
		return err
	}

	repo, err := postgres.New(cfg.DatabaseDSN, l)
	if err != nil {
		return err
	}
	defer multierr.AppendInvoke(&err, multierr.Close(repo))

	s := service.New(repo, l, service.WithJWTSecret(cfg.JWTSecret), service.WithMasterKey(cfg.MasterKey))

	grpcSrv, err := grpc.New(s, l, grpc.WithGRPCAddress(cfg.GRPCAddr), grpc.WithJWTSecret(cfg.JWTSecret))
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcSrv.Run(); err != nil {
			errCh <- err
		}
	}()

	var runErr error

	select {
	case <-ctx.Done():
		l.Info().Msg("Получен сигнал завершения")
	case err := <-errCh:
		runErr = err
		cancel()
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	errChShutdown := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcSrv.Shutdown(shutdownCtx); err != nil {
			errChShutdown <- err
		}
	}()

	wg.Wait()
	close(errChShutdown)

	for e := range errChShutdown {
		runErr = multierr.Append(runErr, e)
	}

	return runErr
}
