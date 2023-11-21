package config

import (
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

type Config struct {
	App *App `json:"app"`

	Api *Api                   `json:"api"`
	Jwt *Jwt                   `json:"jwt"`
	Pg  *Postgres              `json:"postgres"`
	Rpc map[string]*RPCDetails `json:"rpc"`
}

type App struct {
	LogLevel        int        `json:"logLevel"`
	LogColorEnabled bool       `json:"logColorEnabled"`
	InstanceName    string     `json:"instanceName"`
	InstanceLabel   string     `json:"instanceLabel"`
	Telemetry       *Telemetry `json:"telemetry"`
}

type Telemetry struct {
	Enabled      bool   `json:"Enabled"`
	NewrelicName string `json:"nrName"`
	NewrelicKey  string `json:"nrKey"`
	DebugEnabled bool   `json:"debugEnabled"`
}

type Api struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Jwt struct {
	SecretKey  string `json:"secretKey"`
	Expiration int    `json:"expiration"`
}

type Postgres struct {
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	User     string `json:"user"`
	Db       string `json:"db"`
	Password string `json:"password"`
	CertPath string `json:"certPath"`
	MaxConns int32  `json:"maxConns"`
}

type RPCDetails struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

func New(configFileName string) (*Config, error) {
	cfg := &Config{}
	configFileName, _ = filepath.Abs(configFileName)
	log.Printf("Loading config: %v", configFileName)

	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Error().Err(err).Msg("Config read error")
		return nil, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&cfg); err != nil {
		log.Error().Err(err).Msg("Config unmarshal error")
		return nil, err
	}
	return cfg, nil
}
