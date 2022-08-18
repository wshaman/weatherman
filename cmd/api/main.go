package main

import (
	"context"
	"io/ioutil"
	"log"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"github.com/arskydev/weatherman/pkg/db"
	"github.com/arskydev/weatherman/pkg/repository"
	"github.com/arskydev/weatherman/pkg/server"
	"github.com/arskydev/weatherman/pkg/service"
	"github.com/arskydev/weatherman/pkg/web/handlers"
)

type AppConfig struct {
	CONF_PATH     string `yaml:"CONF_PATH"`
	PASS_ENV_NAME string `yaml:"PASS_ENV_NAME"`
	APP_PORT      string `yaml:"APP_PORT"`
}

const (
	APP_CONF_PATH = "config/app_config.yaml"
)

func main() {
	appConfig, err := NewAppConfig(APP_CONF_PATH)
	if err != nil {
		log.Fatal("Error while gathering app config", err)
	}

	pgCfg, err := db.NewPGConfig(appConfig.CONF_PATH, appConfig.PASS_ENV_NAME)

	if err != nil {
		log.Fatal("Error while initiating db config:", err)
	}

	pgDB, err := db.NewConnect(pgCfg)

	if err != nil {
		log.Fatal("Error while connecting to DB:", err)
	}

	repo, err := repository.NewRepository(pgDB)
	if err != nil {
		log.Fatal("Error while initiating repository:", err)
	}

	service := service.NewService(repo)
	handler := handlers.NewHandler(service)
	server := server.NewServer(handler.InitRoutes())

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := server.Run(ctx, appConfig.APP_PORT); err != nil {
		log.Printf("Error raised while server run:\n%s", err.Error())
	}

	log.Println("Closing db connect...")
	if err := pgDB.Close(); err != nil {
		log.Fatalf("Error raised while closing db connect:\n%s", err.Error())
	}
	log.Println("Closing db connect... Done")

}

func NewAppConfig(confPath string) (*AppConfig, error) {
	yfile, err := ioutil.ReadFile(confPath)

	if err != nil {
		return nil, err
	}

	conf := &AppConfig{}
	err = yaml.Unmarshal(yfile, conf)

	if err != nil {
		return nil, err
	}

	return conf, nil
}
