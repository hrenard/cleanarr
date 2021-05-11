package internal

import (
	"fmt"
	"sync"

	units "github.com/docker/go-units"
)

type Sonarr struct {
	Servarr
}

func NewSonarr(config ServarrConfig) Sonarr {
	return Sonarr{
		Servarr: NewServarr(config, "sonarr"),
	}
}

func (s Sonarr) Series() []Serie {
	var series []Serie
	resp, err := s.Request().
		SetResult(&series).
		Get("/series")

	s.handleError(resp, err)

	return series
}

func (s Sonarr) EpisodeFile(fileId int) *EpisodeFile {
	var file EpisodeFile
	resp, err := s.Request().
		SetResult(&file).
		Get(fmt.Sprintf("/episodefile/%d", fileId))

	s.handleError(resp, err)

	return &file
}

func (s Sonarr) SerieEpisodes(serie Serie, c chan []Episode) {
	var episodes []Episode
	resp, err := s.Request().
		SetQueryParam("seriesId", fmt.Sprintf("%d", serie.Id)).
		SetResult(&episodes).
		Get("/episode")

	s.handleError(resp, err)

	for i, episode := range episodes {
		episodes[i].Serie = &serie
		if episode.HasFile {
			episodes[i].File = s.EpisodeFile(episode.FileId)
		}
	}
	c <- episodes
}

func (s Sonarr) Episodes() EpisodeList {
	series := s.Series()

	c := make(chan []Episode)
	for _, serie := range series {
		go s.SerieEpisodes(serie, c)
	}

	var episodes []Episode
	for i := 0; i < len(series); i++ {
		episodes = append(episodes, <-c...)
	}

	return episodes
}

func (s Sonarr) DeleteEpisodes(episodes EpisodeList) {
	for _, episode := range episodes {
		resp, err := s.Request().
			Delete(fmt.Sprintf("/episodefile/%d", episode.FileId))

		s.handleError(resp, err)

		s.log.Debugf("S%dE%d of %s deleted", episode.SeasonNumber, episode.EpisodeNumber, episode.Serie.Title)
	}
	s.log.Infof("%d episodes deleted, %s freed", len(episodes), units.BytesSize(float64(episodes.Size())))
}

func (s Sonarr) UnmonitorEpisodes(episodes EpisodeList) {
	body := episodes.MonitorUpdate(false)

	resp, err := s.Request().
		SetBody(body).
		Put("/episode/monitor")

	s.handleError(resp, err)

	for _, episode := range episodes {
		s.log.Debugf("S%dE%d of %s unmonitored", episode.SeasonNumber, episode.EpisodeNumber, episode.Serie.Title)
	}

	s.log.Infof(" %d epsiodes unmonitored", len(episodes))
}

func (s Sonarr) Tick(wg *sync.WaitGroup) {
	episodes := s.Episodes().
		WithFile().
		ByFileDate()

	var garbage EpisodeList
	totalSize := episodes.Size()

	for _, episode := range episodes {
		if episode.Expired(s.maxDays) || (s.maxBytes > 0 && totalSize-garbage.Size() > s.maxBytes) {
			garbage = append(garbage, episode)
		}
	}

	if len(garbage) > 0 {
		for _, episode := range garbage {
			s.log.Printf("%s S%dE%d", episode.Serie.Title, episode.SeasonNumber, episode.EpisodeNumber)
		}
		s.UnmonitorEpisodes(garbage)
		s.DeleteEpisodes(garbage)
	}

	wg.Done()
}
