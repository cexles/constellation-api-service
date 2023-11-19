package main

import (
	"api-service/internal/api"
	"api-service/internal/api/handler"
	"api-service/internal/config"
	"api-service/internal/repository"
	"api-service/internal/service"
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.New("config.json")
	if err != nil {
		log.Panic().Msgf("Couldn't parse config \n %v ", err)
	}

	initGlobalLogger(cfg.App)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pgCfg, err := config.NewPgConfig(cfg.Pg)
	if err != nil {
		log.Panic().Err(err).Msg("pg cfg")
	}

	pg, err := initPg(ctx, pgCfg)
	if err != nil {
		log.Panic().Err(err).Msg("pg")
	}

	defer func() {
		log.Info().Msg("Shutting down postgres...")
		pg.Close()

	}()

	userRepo := repository.NewUser(pg)
	authService := service.NewAuth(userRepo)
	authHandler := handler.NewApi(authService)
	app := api.NewFiber(ctx, cfg.Jwt, authHandler)

	go func() {
		<-ctx.Done()
		cancel()
		err = app.Shutdown()
		if err != nil {
			log.Error().Msgf("Can't stop server gracefully %v", err)
		}
		log.Info().Msg("Api graceful shutdown, exiting in few seconds...")
	}()

	go func() {
		err = app.Listen(fmt.Sprintf("%s:%d", cfg.Api.Host, cfg.Api.Port))
		if err != nil {
			log.Panic().Msgf("Main / app.Server().Serve(listener) - err", err)
			return
		}
	}()

	log.Info().Msgf("Server started on %s:%d", cfg.Api.Host, cfg.Api.Port)
	<-ctx.Done()
}

func initGlobalLogger(cfg *config.App) {
	devWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat, NoColor: !cfg.LogColorEnabled}
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.InterfaceMarshalFunc = json.Marshal

	log.Logger = log.With().Timestamp().Caller().Str("service", cfg.InstanceName).Str("service_label", cfg.InstanceLabel).Stack().Logger().Output(devWriter)

	logLevel := zerolog.Level(cfg.LogLevel)
	log.Info().Msgf("Log level: %s", logLevel)

	zerolog.SetGlobalLevel(logLevel)
}

func initPg(ctx context.Context, cfg *pgxpool.Config) (*pgxpool.Pool, error) {
	pgPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	t := time.Now().UTC()
	err = pgPool.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Fields(map[string]any{
			"host": cfg.ConnConfig.Host,
			"port": cfg.ConnConfig.Host,
		}).Msg("pg ping error")
		return nil, err
	}
	log.Info().Fields(map[string]any{
		"host":   cfg.ConnConfig.Host,
		"port":   cfg.ConnConfig.Port,
		"db":     cfg.ConnConfig.Database,
		"rtt_ms": time.Since(t).Milliseconds(),
	}).Msg("connected to postgres")

	return pgPool, nil
}
