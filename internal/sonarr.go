package internal

import "sync"

type Sonarr struct {
}

func NewSonarr(config ServarrConfig) Sonarr {
	return Sonarr{}
}

func (s Sonarr) Tick(wg *sync.WaitGroup) {
	wg.Done()
}
