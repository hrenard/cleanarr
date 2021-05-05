package main

import (
	"os"
	"sync"
	"time"

	"github.com/hrenard/cleanarr/radarr"
	"github.com/hrenard/cleanarr/servarr"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CleanarrConfig struct {
	Interval int
	Radarr   []servarr.Config
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

	servarrList := make([]servarr.Servarr, len(config.Radarr))
	for _, radarrConf := range config.Radarr {
		servarrList[0] = radarr.New(&radarrConf)
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
