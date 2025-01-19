package internal

import (
	"context"
	"fmt"

	"github.com/alitto/pond/v2"
	units "github.com/docker/go-units"
	"golift.io/starr/sonarr"
)

type Sonarr struct {
	*sonarr.Sonarr
	*Servarr
}

func NewSonarr(config ServarrConfig) Sonarr {
	servarr, starrConf := ParseConfig(config)
	servarr.log = servarr.log.WithField("app", "sonarr")
	return Sonarr{
		Servarr: &servarr,
		Sonarr:  sonarr.New(starrConf),
	}
}

func (s Sonarr) RefreshTags(ctx context.Context) error {
	tags, err := s.GetTagsContext(ctx)
	if err == nil {
		s.refreshTags(tags)
	}
	return err
}

func (s Sonarr) Fetch(ctx context.Context) (Cleanables, error) {
	series, err := s.GetSeriesContext(ctx, 0)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	pool := pond.NewResultPool[Cleanables](10)
	tasks := make([]pond.Result[Cleanables], len(series))
	for i, serie := range series {
		serie := serie
		tasks[i] = pool.SubmitErr(func() (Cleanables, error) {
			episodes, err := s.GetSeriesEpisodesContext(ctx, serie.ID)
			if err != nil {
				return nil, err
			}
			files, err := s.GetSeriesEpisodeFilesContext(ctx, serie.ID)
			if err != nil {
				return nil, err
			}
			filesMap := map[int64]*sonarr.EpisodeFile{}
			for _, file := range files {
				filesMap[file.ID] = file
			}
			cleanables := Cleanables{}
			for _, episode := range episodes {
				if episode.HasFile {
					file, ok := filesMap[episode.EpisodeFileID]
					if !ok {
						return nil, fmt.Errorf("file not found for episode %d", episode.ID)
					}
					cleanables = append(cleanables, Episode{
						Series:      serie,
						Episode:     episode,
						EpisodeFile: file,
					})
				}
			}
			return cleanables, nil
		})
	}
	cleanables := make(Cleanables, 0)
	for _, task := range tasks {
		c, err := task.Wait()
		if err != nil {
			pool.Stop()
			return nil, err
		}
		cleanables = append(cleanables, c...)
	}
	return cleanables, nil
}

func (s Sonarr) Clean(ctx context.Context, cleanables Cleanables, dryRun bool) error {
	episodes := []Episode{}
	ids := []int64{}
	for _, cleanable := range cleanables {
		episode, ok := cleanable.(Episode)
		if !ok {
			return fmt.Errorf("cleanable is not of type Episode")
		}
		episodes = append(episodes, episode)
		ids = append(ids, episode.Episode.ID)
	}
	if !dryRun {
		_, err := s.MonitorEpisodeContext(ctx, ids, false)
		if err != nil {
			return err
		}
	}
	for _, episode := range episodes {
		s.log.Debugf("S%dE%d of %s unmonitored", episode.Episode.SeasonNumber, episode.EpisodeNumber, episode.Series.Title)
	}
	s.log.Infof("%d epsiodes unmonitored", len(cleanables))

	for _, episode := range episodes {
		if !dryRun {
			err := s.DeleteEpisodeFileContext(ctx, episode.EpisodeFileID)
			if err != nil {
				return err
			}
		}
		s.log.Debugf("S%dE%d of %s deleted", episode.Episode.SeasonNumber, episode.EpisodeNumber, episode.Series.Title)
	}
	s.log.Infof("%d episodes deleted, %s freed", len(cleanables), units.BytesSize(float64(cleanables.Size())))
	return nil
}
