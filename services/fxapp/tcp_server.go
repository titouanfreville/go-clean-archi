package fxapp

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"go-clean-archi/servers"
)

type server struct {
	lifecycle  fx.Lifecycle
	shutdowner fx.Shutdowner
	logger     *zap.Logger
}

func NewTCPServer(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, logger *zap.Logger) *server {
	return &server{lifecycle: lifecycle, shutdowner: shutdowner, logger: logger}
}

func (server server) Run(name string, transport servers.TCP) {
	logger := server.logger.
		Named("Lifecycle").
		With(zap.String("address", transport.GetAddress()))

	server.lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				go func() {
					if err := transport.ListenAndServe(); err != nil {
						if !ShuttingDown {
							logger.Error("🧨💥 "+name+" closed unexpectedly", zap.Error(err))
						}

						if err = server.shutdowner.Shutdown(); err != nil {
							logger.Error("🧨💥 Unable to shutdown properly "+name, zap.Error(err))
						}
					}
				}()

				logger.Info("📢 " + name + " started")

				return nil
			},

			OnStop: func(context.Context) error {
				if err := transport.Shutdown(); err != nil {
					return err
				}

				logger.Info("📢 " + name + " closed")

				return nil
			},
		},
	)
}
