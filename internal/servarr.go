package internal

import (
	"encoding/json"
	"sync"

	units "github.com/docker/go-units"
	resty "github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

type ServarrConfig struct {
	Name     string
	HostPath string
	ApiKey   string
	MaxDays  *int
	MaxSize  *string
	MaxFiles *int
}

type IServarr interface {
	Tick(*sync.WaitGroup)
}

type Servarr struct {
	client   *resty.Client
	log      *log.Entry
	maxDays  int
	maxBytes int
}

type jsonError struct {
	Error string `json:"error"`
}

func NewServarr(config ServarrConfig, app string) Servarr {
	servarr := Servarr{
		client: resty.New().
			SetHostURL(config.HostPath+"/api/v3").
			SetQueryParam("apikey", config.ApiKey),
		log: log.WithFields(log.Fields{
			"app":  app,
			"name": config.Name,
		}),
	}

	if config.MaxDays == nil && config.MaxSize == nil {
		servarr.log.Fatal("No constraints, maxDays or maxSize required")
	}

	if config.MaxDays != nil {
		servarr.maxDays = *config.MaxDays
	}

	if config.MaxSize != nil {
		maxBytes, err := units.FromHumanSize(*config.MaxSize)
		if err != nil {
			servarr.log.Fatalf("Failed to parse maxSize: %s", err)
		}
		servarr.maxBytes = int(maxBytes)
	}

	return servarr
}

func (s Servarr) Request() *resty.Request {
	return s.client.R()
}

func (s Servarr) handleError(resp *resty.Response, err error) {
	if err != nil {
		s.log.Errorf("request failed: %s", err)
	}

	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		var jsonErr jsonError
		errLog := s.log.WithField("http_status", resp.StatusCode())
		if err := json.Unmarshal(resp.Body(), &jsonErr); err != nil {
			errLog.Errorf("request error: %s", resp.Body())
		} else {
			errLog.Errorf("request error: %s", jsonErr.Error)
		}
	}
}
