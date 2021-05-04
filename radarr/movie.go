package radarr

import "time"

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

func (m Movie) HasFile() bool {
	return m.File != nil
}

func (m Movie) Expired(maxDays int) bool {
	dateLimit := m.File.DateAdded.AddDate(0, 0, maxDays)
	return time.Now().After(dateLimit)
}
