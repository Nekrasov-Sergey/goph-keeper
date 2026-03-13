package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"

	"github.com/Nekrasov-Sergey/goph-keeper/internal/client"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/client/grpc"
	"github.com/Nekrasov-Sergey/goph-keeper/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Клиент завершился с ошибкой")
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cfg, err := config.NewClientConfig()
	if err != nil {
		return err
	}

	grpcClient, err := grpc.New(cfg.GRPCAddr)
	if err != nil {
		return err
	}
	defer multierr.AppendInvoke(&err, multierr.Close(grpcClient))

	cli := client.New(cfg, grpcClient.Client)

	cli.AuthMenu(ctx)
	fmt.Println()
	fmt.Println("Программа завершена")

	return nil
}
