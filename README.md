# Cleanarr

A small utility tasked to automatically clean radarr and sonarr files over time.

The idea is to have a "rolling" media library where old media is deleted either because the disk quota is exceeded or because it has been there too long.

## Usage

### Configuration

```yml
#config.yml
interval: 1 # optional, check every minutes
dryRun: false # optional, do not perform actions

radarr:
  - name: "radarr4k" # requried
    hostPath: "https://radarr.mydomain.com" # required
    apiKey: "xxxxxxxxxxx" # required
    maxDays: 90 # optional if maxSize
    maxSize: "2TB" # optional if maxDays
    includeTags: # optional, includes all by default
      - rolling
    excludeTags: [] # optional, excludes nothing by default

sonarr:
  - name: "sonarr" # requried
    hostPath: "https://sonarr.mydomain.com" # required
    apiKey: "xxxxxxxxxxx" # required
    maxDays: 90 # optional if maxSize
    maxSize: "2TB" # optional if maxDays
    includeTags: [] # optional, includes all by default
    excludeTags: # optional, excludes nothing by default
      - keep 
```

### Docker

```shell
$ docker run -v $PWD/config.yml:/config.yml ghcr.io/hrenard/cleanarr
```

## Roadmap

- [ ] Radarr
  - [x] Days policy
  - [x] Size policy
  - [x] Tag filter
  - [ ] Quantity policy
- [ ] Sonarr
  - [x] Days policy
  - [x] Size policy
  - [x] Tag filter
  - [ ] Quantity policy
  - [ ] Unmonitor season when all episodes are unmonitored ?
  - [ ] Remove serie when all episodes are unmonitored ?
