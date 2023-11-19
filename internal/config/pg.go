package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

func NewPgConfig(cfg *Postgres) (*pgxpool.Config, error) {
	pgCfg, err := pgxpool.ParseConfig(fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Db))
	if err != nil {
		log.Error().Err(err).Fields(map[string]any{
			"host": cfg.Host,
			"port": cfg.Port,
		}).Msg("pg config error")
		return nil, err
	}

	pgCfg.MaxConns = cfg.MaxConns
	if len(cfg.CertPath) > 0 {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(cfg.CertPath)
		if err != nil {
			panic(err)
		}

		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			panic("Failed to append PEM.")
		}

		pgCfg.ConnConfig.TLSConfig = &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		}
	}

	return pgCfg, nil
}
