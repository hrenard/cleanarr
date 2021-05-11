package main

import (
	"os"
	"sync"
	"time"

	"github.com/hrenard/cleanarr/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CleanarrConfig struct {
	Interval int
	Radarr   []internal.ServarrConfig
	Sonarr   []internal.ServarrConfig
}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	// log.SetLevel(log.DebugLevel)

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

	servarrList := make([]internal.IServarr, len(config.Radarr)+len(config.Sonarr))

	for i, radarrConf := range config.Radarr {
		servarrList[i] = internal.NewRadarr(radarrConf)
	}

	for i, sonarrConf := range config.Sonarr {
		servarrList[len(config.Radarr)+i] = internal.NewSonarr(sonarrConf)
	}

	log.Infof("Cleanarr is running")
	for {
		var wg sync.WaitGroup
		for _, s := range servarrList {
			wg.Add(1)
			go s.Tick(&wg)
		}
		wg.Wait()
		time.Sleep(time.Minute * time.Duration(config.Interval))
	}
}
