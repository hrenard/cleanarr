# Cleanarr

A small utility tasked to automatically clean radarr and sonarr files over time.

## Usage

### Configuration

```yml
#config.yml

radarr:
  - name: "radarr4k" # requried
    hostPath: "https://radarr.mydomain.com" # required
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
  - [ ] Quantity polixy
- [ ] Sonarr
  - [ ] Days policy
  - [ ] Size policy
  - [ ] Quantity polixy
