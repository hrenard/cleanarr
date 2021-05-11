# Cleanarr

A small utility tasked to automatically clean radarr and sonarr files over time.

## Usage

### Configuration

```yml
#config.yml
interval: 1 # optional, check every minutes

radarr:
  - name: "radarr4k" # requried
    hostPath: "https://radarr.mydomain.com" # required
    apiKey: "xxxxxxxxxxx" # required
    maxDays: 90 # optional if maxSize
    maxSize: "2TB" # optional if maxDays

sonarr:
  - name: "sonarr" # requried
    hostPath: "https://sonarr.mydomain.com" # required
    apiKey: "xxxxxxxxxxx" # required
    maxDays: 90 # optional if maxSize
    maxSize: "2TB" # optional if maxDays
```

### Docker

```shell
$ docker run -v $PWD/config.yml:/config.yml ghcr.io/hrenard/cleanarr
```

## Roadmap

- [ ] Radarr
  - [x] Days policy
  - [x] Size policy
  - [ ] Quantity policy
- [ ] Sonarr
  - [x] Days policy
  - [x] Size policy
  - [ ] Quantity policy
  - [ ] Unmonitor season when all episodes are unmonitored ?
  - [ ] Remove serie when all episodes are unmonitored ?
