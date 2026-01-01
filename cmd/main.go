package main

import (
	"log"
	"os"

	"github.com/dormitory-life/gateway/internal/config"
	"github.com/dormitory-life/gateway/internal/logger"
	"github.com/dormitory-life/gateway/internal/server"
)

func main() {
	configPath := os.Args[1]
	cfg, err := config.ParseConfig(configPath)
	if err != nil {
		panic(err)
	}

	log.Println("CONFIG: ", cfg)

	logger, err := logger.New(cfg)
	if err != nil {
		panic(err)
	}

	s := server.New(server.ServerConfig{
		Config:         cfg.Server,
		AuthServiceUrl: cfg.AuthService.Url,
		CoreServiceUrl: cfg.CoreService.Url,
		Logger:         logger,
		JWTSecret:      cfg.JWT.Secret,
	})

	panic(s.Start())
}
