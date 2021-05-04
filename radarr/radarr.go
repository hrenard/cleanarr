package radarr

import (
	"encoding/json"
	"fmt"
	"sync"

	units "github.com/docker/go-units"
	resty "github.com/go-resty/resty/v2"
	"github.com/hrenard/cleanarr/servarr"
	log "github.com/sirupsen/logrus"
)

type Radarr struct {
	client   *resty.Client
	log      *log.Entry
	maxDays  int
	maxBytes int
}

type jsonError struct {
	Error string `json:"error"`
}

func New(config *servarr.Config) Radarr {
	radarr := Radarr{
		client: resty.New().
			SetHostURL(config.HostPath+"/api/v3").
			SetQueryParam("apikey", config.ApiKey),
		log: log.WithFields(log.Fields{
			"app":  "radarr",
			"name": config.Name,
		}),
	}

	if config.MaxDays == nil && config.MaxSize == nil {
		radarr.log.Fatal("No constraints, maxDays or maxSize required")
	}

	if config.MaxDays != nil {
		radarr.maxDays = *config.MaxDays
	}

	if config.MaxSize != nil {
		maxBytes, err := units.FromHumanSize(*config.MaxSize)
		if err != nil {
			radarr.log.Fatalf("Failed to parse maxSize: %s", err)
		}
		radarr.maxBytes = int(maxBytes)
	}

	return radarr
}

func (r Radarr) Request() *resty.Request {
	return r.client.R()
}

func (r Radarr) handleError(resp *resty.Response, err error) {
	if err != nil {
		r.log.Errorf("request failed: %s", err)
	}

	if resp.StatusCode() != 200 {
		var jsonErr jsonError
		errLog := r.log.WithField("http_status", resp.StatusCode())
		if err := json.Unmarshal(resp.Body(), &jsonErr); err != nil {
			errLog.Errorf("request error: %s", resp.Body())
		} else {
			errLog.Errorf("request error: %s", jsonErr.Error)
		}
	}
}

func (r Radarr) Movies() []Movie {
	resp, err := r.Request().Get("/movie")
	r.handleError(resp, err)

	var movies []Movie
	if err := json.Unmarshal(resp.Body(), &movies); err != nil {
		r.log.Error(err)
	}

	return movies
}

func (r Radarr) DeleteMovie(movie *Movie) {
	resp, err := r.Request().
		SetQueryParam("deleteFiles", "true").
		Delete(fmt.Sprintf("/movie/%d", movie.Id))
	r.handleError(resp, err)
	r.log.Printf("%s deleted", movie.Title)
}

func (r Radarr) PurgeExpiredMovies() {
	movies := r.Movies()
	for _, movie := range movies {
		if movie.HasFile() {
			if r.maxDays > 0 {
				if movie.Expired(r.maxDays) {
					r.log.Debugf("%s expired", movie.Title)
					r.DeleteMovie(&movie)
				}
			}
		}
	}
}

func (r Radarr) PurgeFirstInMovies() {
	movies := r.Movies()
	var totalBytes int
	var oldestFile *Movie
	for _, movie := range movies {
		if movie.HasFile() {
			totalBytes += movie.File.Size
			if oldestFile == nil || oldestFile.File.DateAdded.After(*movie.File.DateAdded) {
				m := movie
				oldestFile = &m
			}
		}
	}
	if totalBytes > r.maxBytes {
		r.log.Debugf("%s in excess", units.HumanSize(float64(totalBytes-r.maxBytes)))
		r.DeleteMovie(oldestFile)
		if totalBytes-oldestFile.File.Size > r.maxBytes {
			r.PurgeFirstInMovies()
		}
	}
}

func (r Radarr) Tick(wg *sync.WaitGroup) {
	r.PurgeExpiredMovies()
	r.PurgeFirstInMovies()
	wg.Done()
}
