package main

import (
	"api-service/internal/api"
	"api-service/internal/api/handler"
	"api-service/internal/config"
	"api-service/internal/repository"
	"api-service/internal/service"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"math/big"
	"os"
	"os/signal"
	"sync"
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
		log.Panic().Msgf("Couldn't parse PG config \n $v", err)
	}

	pg, err := initPg(ctx, pgCfg)
	if err != nil {
		log.Panic().Msgf("PG connection error \n $v", err)
	}

	defer func() {
		log.Info().Msg("Shutting down postgres...")
		pg.Close()

	}()

	ethClients := make(map[string]*ethclient.Client, len(cfg.Rpc))
	mu := &sync.Mutex{}
	for chain, r := range cfg.Rpc {
		if !r.Enabled {
			log.Info().Msgf("Chain disabled: %s", chain)
			continue
		}

		go func(chain string, r *config.RPCDetails, mu *sync.Mutex) {
			var client *ethclient.Client
			var chainId *big.Int

			for {
				c, err := ethclient.DialContext(ctx, r.URL)
				if err == nil {
					chainId, err = c.ChainID(ctx)
				}
				if err == nil {
					client = c
					log.Info().Fields(map[string]any{
						"rpc":      r.URL,
						"chain_id": chainId,
					}).Msg("Connected to RPC")
					break
				}
				log.Warn().Msgf("Error connecting: %s", r.URL)
			}
			mu.Lock()
			ethClients[chain] = client
			mu.Unlock()
		}(chain, r, mu)
	}

	defer func() {
		log.Info().Msg("Shutting down RPC clients")
		for _, v := range ethClients {
			v.Close()
		}
	}()

	userRepo := repository.NewUser(pg)
	chainRepo := repository.NewChain(pg)
	tokenRepo := repository.NewToken(pg)

	authService := service.NewAuth(userRepo, *cfg.Jwt)
	chainService := service.NewChain(chainRepo, cfg.Rpc)
	tokenSevice := service.NewToken(tokenRepo)
	balanceService := service.NewBalance(ethClients, *chainService, *tokenSevice)

	authHandler := handler.NewAuthApi(authService)
	userHandler := handler.NewUserApi(chainService, balanceService)

	app := api.NewFiber(ctx, cfg.Jwt, authHandler, userHandler)

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
		}).Msg("PG Ping error")
		return nil, err
	}
	log.Info().Fields(map[string]any{
		"host":   cfg.ConnConfig.Host,
		"port":   cfg.ConnConfig.Port,
		"db":     cfg.ConnConfig.Database,
		"rtt_ms": time.Since(t).Milliseconds(),
	}).Msg("Connected to PG")

	return pgPool, nil
}
