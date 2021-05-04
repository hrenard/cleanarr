package servarr

import (
	"sync"
)

type Config struct {
	Name     string
	HostPath string
	ApiKey   string
	MaxDays  *int
	MaxSize  *string
}

type Servarr interface {
	Tick(*sync.WaitGroup)
}
