package internal

import (
	"context"
	"slices"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type Cleanable interface {
	Size() int
	Date() time.Time
	GetTags() []int
}

type Cleanables []Cleanable

func (cleanables Cleanables) Size() int {
	var size int
	for _, cleanable := range cleanables {
		size += cleanable.Size()
	}
	return size
}

type Provider interface {
	RefreshTags(context.Context) error
	Fetch(context.Context) (Cleanables, error)
	Clean(context.Context, Cleanables, bool) error
	MaxDays() int
	MaxBytes() int
	IncludeTagIDs() []int
	ExcludeTagIDs() []int
	Log() *log.Entry
}

func CronProcess(ctx context.Context, wg *sync.WaitGroup, provider Provider, interval int, dryRun bool) {
	defer wg.Done()
	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			Process(ctx, provider, dryRun)
		}
	}
}

func Process(ctx context.Context, provider Provider, dryRun bool) {
	provider.Log().Info("Processing")

	err := provider.RefreshTags(ctx)
	if err != nil {
		provider.Log().Error(err)
		return
	}

	cleanables, err := provider.Fetch(ctx)
	if err != nil {
		provider.Log().Error(err)
		return
	}
	sort.Slice(cleanables, func(i, j int) bool {
		return cleanables[i].Date().Before(cleanables[j].Date())
	})

	var garbage Cleanables
	for _, cleanable := range cleanables {
		if ShouldSkip(provider, cleanable) {
			continue
		}
		if IsExpired(cleanable, provider.MaxDays()) || IsExceedingQuota(provider, cleanables, garbage) {
			garbage = append(garbage, cleanable)
		}
	}

	if len(garbage) > 0 {
		err = provider.Clean(ctx, garbage, dryRun)
		if err != nil {
			provider.Log().Error(err)
		}
	}

	provider.Log().Info("Done")
}

func IsExpired(cleanable Cleanable, maxDays int) bool {
	if maxDays == 0 {
		return false
	}
	dateLimit := cleanable.Date().AddDate(0, 0, maxDays)
	return time.Now().After(dateLimit)
}

func IsExceedingQuota(provider Provider, cleanables, garbage Cleanables) bool {
	return provider.MaxBytes() > 0 && cleanables.Size()-garbage.Size() > provider.MaxBytes()
}

func ShouldSkip(provider Provider, cleanable Cleanable) bool {
	for _, tagID := range provider.ExcludeTagIDs() {
		if slices.Contains(cleanable.GetTags(), tagID) {
			return true
		}
	}
	for _, tagID := range provider.IncludeTagIDs() {
		if slices.Contains(cleanable.GetTags(), tagID) {
			return false
		}
	}
	return len(provider.IncludeTagIDs()) > 0
}
