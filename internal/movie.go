package internal

import (
	"sort"
	"time"
)

type Movie struct {
	Id    int        `json:"id"`
	Title string     `json:"title"`
	Added *time.Time `json:"added"`
	File  *MovieFile `json:"movieFile"`
}

type MovieFile struct {
	Size      int        `json:"size"`
	DateAdded *time.Time `json:"dateAdded"`
}

type MovieList []Movie

func (m Movie) HasFile() bool {
	return m.File != nil
}

func (m Movie) Expired(maxDays int) bool {
	if maxDays == 0 {
		return false
	}

	dateLimit := m.File.DateAdded.AddDate(0, 0, maxDays)
	return time.Now().After(dateLimit)
}

func (movies MovieList) Size() int {
	var size int
	for _, movie := range movies {
		if movie.HasFile() {
			size += movie.File.Size
		}
	}
	return size
}

func (movies MovieList) WithFile() MovieList {
	var filterd MovieList
	for _, movie := range movies {
		if movie.HasFile() {
			filterd = append(filterd, movie)
		}
	}
	return filterd
}

func (movies MovieList) ByFileDate() MovieList {
	sort.Slice(movies, func(i, j int) bool {
		return movies[i].File.DateAdded.Before(*movies[j].File.DateAdded)
	})
	return movies
}
