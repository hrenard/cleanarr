package internal

import (
	"os"
	"slices"
	"strings"

	units "github.com/docker/go-units"
	log "github.com/sirupsen/logrus"
	"golift.io/starr"
)

type ServarrConfig struct {
	Name        string
	HostPath    string
	ApiKey      string
	ApiKeyFile  string
	MaxDays     *int
	MaxSize     *string
	MaxFiles    *int
	IncludeTags []string
	ExcludeTags []string
}

type Servarr struct {
	log            *log.Entry
	maxDays        int
	maxBytes       int
	includeTags    []string
	excludeTags    []string
	includeTagIDs  []int
	excludeTagsIDs []int
}

func ParseConfig(config ServarrConfig) (Servarr, *starr.Config) {
	servarr := Servarr{
		log: log.WithFields(log.Fields{
			"name": config.Name,
		}),
		includeTags: config.IncludeTags,
		excludeTags: config.ExcludeTags,
	}

	if config.MaxDays == nil && config.MaxSize == nil {
		servarr.log.Fatal("No constraints, maxDays or maxSize required")
	}

	if config.MaxDays != nil {
		servarr.maxDays = *config.MaxDays
	}

	if config.MaxSize != nil {
		maxBytes, err := units.FromHumanSize(*config.MaxSize)
		if err != nil {
			servarr.log.Fatalf("Failed to parse maxSize: %s", err)
		}
		servarr.maxBytes = int(maxBytes)
	}

	apiKey := config.ApiKey
	if apiKey == "" {
		if config.ApiKeyFile == "" {
			servarr.log.Fatal("apiKey or apiKeyFile must be set")
		}
		rawKey, err := os.ReadFile(config.ApiKeyFile)
		if err != nil {
			servarr.log.Fatalf("Failed to read apiKeyFile %s: %s", config.ApiKeyFile, err)
		}
		apiKey = strings.TrimSpace(string(rawKey))
	}

	return servarr, starr.New(apiKey, config.HostPath, 0)
}

func (s *Servarr) Log() *log.Entry {
	return s.log
}

func (s *Servarr) MaxDays() int {
	return s.maxDays
}

func (s *Servarr) MaxBytes() int {
	return s.maxBytes
}

func (s *Servarr) IncludeTagIDs() []int {
	return s.includeTagIDs
}

func (s *Servarr) ExcludeTagIDs() []int {
	return s.excludeTagsIDs
}

func (s *Servarr) refreshTags(tags []*starr.Tag) {
	s.includeTagIDs = make([]int, len(s.includeTags))
	s.excludeTagsIDs = make([]int, len(s.excludeTags))
	for _, tag := range tags {
		if i := slices.Index(s.includeTags, tag.Label); i > -1 {
			s.includeTagIDs[i] = tag.ID
		}
		if i := slices.Index(s.excludeTags, tag.Label); i > -1 {
			s.excludeTagsIDs[i] = tag.ID
		}
	}
}
