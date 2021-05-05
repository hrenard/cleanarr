package internal

import (
	"sync"
)

type ServarrConfig struct {
	Name     string
	HostPath string
	ApiKey   string
	MaxDays  *int
	MaxSize  *string
	MaxFiles *int
}

type Servarr interface {
	Tick(*sync.WaitGroup)
}
