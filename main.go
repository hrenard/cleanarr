package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hrenard/cleanarr/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CleanarrConfig struct {
	Interval int
	DryRun   bool
	Radarr   []internal.ServarrConfig
	Sonarr   []internal.ServarrConfig
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.SetDefault("interval", 1)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %s \n", err.Error())
	}

	var config CleanarrConfig
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to parse config: %s \n", err.Error())
	}

	if config.Radarr == nil {
		log.Fatalf("No radarr configured")
	}

	providers := make([]internal.Provider, len(config.Radarr)+len(config.Sonarr))

	for i, radarrConf := range config.Radarr {
		providers[i] = internal.NewRadarr(radarrConf)
	}

	for i, sonarrConf := range config.Sonarr {
		providers[len(config.Radarr)+i] = internal.NewSonarr(sonarrConf)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	for _, provider := range providers {
		wg.Add(1)
		go internal.CronProcess(ctx, &wg, provider, config.Interval, config.DryRun)
	}
	log.Infof("Cleanarr is running")
	wg.Wait()
}
