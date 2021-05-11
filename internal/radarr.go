package internal

import (
	"encoding/json"
	"fmt"
	"sync"

	units "github.com/docker/go-units"
)

type Radarr struct {
	Servarr
}

func NewRadarr(config ServarrConfig) Radarr {
	return Radarr{
		Servarr: NewServarr(config, "radarr"),
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
	resp, err := r.Request().
		SetQueryParam("deleteFiles", "true").
		Delete(fmt.Sprintf("/movie/%d", movie.Id))
	r.handleError(resp, err)
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
