package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/wxvn/grpc-inventory-service/internal/app/grpc"
	"github.com/wxvn/grpc-inventory-service/internal/repository/postgres"
	"github.com/wxvn/grpc-inventory-service/internal/service"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(ctx context.Context, log *slog.Logger, grpcPort int, dns string) *App {
	repos, err := postgres.New(ctx, dns)
	if err != nil {
		panic(err)
	}

	inventoryService := service.New(repos, log)

	grpcApp := grpcapp.New(
		log,
		inventoryService,
		grpcPort,
	)

	return &App{
		GRPCServer: grpcApp,
	}
}
