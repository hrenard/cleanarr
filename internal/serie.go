package internal

import (
	"time"

	"golift.io/starr/sonarr"
)

type Episode struct {
	*sonarr.Series
	*sonarr.Episode
	*sonarr.EpisodeFile
}

func (episode Episode) Size() int {
	return int(episode.EpisodeFile.Size)
}

func (episode Episode) Date() time.Time {
	return episode.EpisodeFile.DateAdded
}

func (episode Episode) GetTags() []int {
	return episode.Series.Tags
}
