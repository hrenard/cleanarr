package internal

import (
	"time"

	"golift.io/starr/radarr"
)

type Movie radarr.Movie

func (movie Movie) Size() int {
	return int(movie.MovieFile.Size)
}

func (movie Movie) Date() time.Time {
	return movie.MovieFile.DateAdded
}

func (movie Movie) GetTags() []int {
	return movie.Tags
}
