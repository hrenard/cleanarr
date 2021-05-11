package internal

import (
	"sort"
	"time"
)

type Serie struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type Episode struct {
	Id            int  `json:"id"`
	FileId        int  `json:"episodeFileId"`
	SeasonNumber  int  `json:"seasonNumber"`
	EpisodeNumber int  `json:"episodeNumber"`
	HasFile       bool `json:"hasFile"`
	Monitored     bool `json:"monitored"`
	Serie         *Serie
	File          *EpisodeFile
}

type EpisodeFile struct {
	Id        int        `json:"id"`
	Size      int        `json:"size"`
	DateAdded *time.Time `json:"dateAdded"`
	// EpisodeId int        `json:"episodeId"`
}

type MonitorUpdate struct {
	Monitored  bool  `json:"monitored"`
	EpisodeIds []int `json:"episodeIds"`
}

type EpisodeList []Episode

func (episode Episode) Expired(maxDays int) bool {
	if maxDays == 0 {
		return false
	}

	dateLimit := episode.File.DateAdded.AddDate(0, 0, maxDays)
	return time.Now().After(dateLimit)
}

func (episodes EpisodeList) Size() int {
	var size int
	for _, episode := range episodes {
		size += episode.File.Size
	}
	return size
}

func (episodes EpisodeList) WithFile() EpisodeList {
	var filtered EpisodeList
	for _, episode := range episodes {
		if episode.HasFile {
			filtered = append(filtered, episode)
		}
	}
	return filtered
}

func (episodes EpisodeList) ByFileDate() EpisodeList {
	sort.Slice(episodes, func(i, j int) bool {
		return episodes[i].File.DateAdded.Before(*episodes[j].File.DateAdded)
	})
	return episodes
}

func (episodes EpisodeList) MonitorUpdate(monitored bool) MonitorUpdate {
	monitorUpdate := MonitorUpdate{
		Monitored:  monitored,
		EpisodeIds: make([]int, len(episodes)),
	}
	for i, episode := range episodes {
		monitorUpdate.EpisodeIds[i] = episode.Id
	}
	return monitorUpdate
}
