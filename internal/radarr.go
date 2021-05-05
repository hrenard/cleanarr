package internal

import (
	"encoding/json"
	"sync"

	units "github.com/docker/go-units"
	resty "github.com/go-resty/resty/v2"
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

func NewRadarr(config ServarrConfig) Radarr {
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

func (r Radarr) Movies() MovieList {
	resp, err := r.Request().Get("/movie")
	r.handleError(resp, err)

	var movies []Movie
	if err := json.Unmarshal(resp.Body(), &movies); err != nil {
		r.log.Error(err)
	}

	return movies
}

func (r Radarr) DeleteMovie(movie Movie) {
	// resp, err := r.Request().
	// 	SetQueryParam("deleteFiles", "true").
	// 	Delete(fmt.Sprintf("/movie/%d", movie.Id))
	// r.handleError(resp, err)
	r.log.Debugf("%s deleted", movie.Title)
}

func (r Radarr) DeleteMovies(movies MovieList) {
	for _, movie := range movies {
		r.DeleteMovie(movie)
	}
	r.log.Printf("%d movies deleted, %s freed", len(movies), units.BytesSize(float64(movies.Size())))
}

func (r Radarr) Tick(wg *sync.WaitGroup) {
	movies := r.Movies().
		WithFile().
		ByFileDate()

	var garbage MovieList
	totalSize := movies.Size()

	for _, movie := range movies {
		if movie.Expired(r.maxDays) || (r.maxBytes > 0 && totalSize-garbage.Size() > r.maxBytes) {
			garbage = append(garbage, movie)
		}
	}

	if len(garbage) > 0 {
		r.DeleteMovies(garbage)
	}

	wg.Done()
}
