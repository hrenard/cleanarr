package internal

import (
	"context"
	"fmt"

	units "github.com/docker/go-units"
	"golift.io/starr/radarr"
)

type Radarr struct {
	*radarr.Radarr
	*Servarr
}

func NewRadarr(config ServarrConfig) Radarr {
	servarr, starrConf := ParseConfig(config)
	servarr.log = servarr.log.WithField("app", "radarr")
	return Radarr{
		Servarr: &servarr,
		Radarr:  radarr.New(starrConf),
	}
}

func (r Radarr) RefreshTags(ctx context.Context) error {
	tags, err := r.GetTagsContext(ctx)
	if err == nil {
		r.refreshTags(tags)
	}
	return err
}

func (r Radarr) Fetch(ctx context.Context) (Cleanables, error) {
	movies, err := r.GetMovieContext(ctx, 0)
	if err != nil {
		return nil, err
	}
	cleanables := make(Cleanables, 0)
	for _, movie := range movies {
		if movie.HasFile {
			cleanables = append(cleanables, Movie(*movie))
		}
	}
	return cleanables, nil
}

func (r Radarr) Clean(ctx context.Context, cleanables Cleanables, dryRun bool) error {
	for _, cleanable := range cleanables {
		movie, ok := cleanable.(Movie)
		if !ok {
			return fmt.Errorf("cleanable is not of type Movie")
		}
		if !dryRun {
			err := r.DeleteMovieContext(ctx, movie.ID, true, false)
			if err != nil {
				return err
			}
		}
		r.log.Debugf("%s deleted", movie.Title)
	}
	r.log.Printf("%d movies deleted, %s freed", len(cleanables), units.BytesSize(float64(cleanables.Size())))
	return nil
}
