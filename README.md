# Uber Zap logger helpers
[![GitHub release][Release img]][Release src] [![Conventional Commits][Conventional commits badge]][Conventional commits src]

## Conf
Features:
- easiest Zap Logger configuration using yaml-config.
- configuration TCP and UDP Zap Logger sinks out of box.
- minimal code for full-featured Zap Logger initialization. 

### Getting started

```go
package main

import (
	"log"
	
	"github.com/nafigator/zapper"
	"go.uber.org/zap"
)

func main() {
	var zl *zap.SugaredLogger
	var err error

	if zl, err = zapper.New(nil, nil); err != nil {
		log.Fatal("Zap logger failure: ", err)
	}
	
	zl.Info("Zap logger: OK")
}
```

### Default configuration:
```yaml
level: info
encoding: console
outputPaths:
  - stdout
errorOutputPaths:
  - stderr
encoderConfig:
  messageKey: message
  levelKey:   level
  timeKey:    time
  callerKey:  line
  levelEncoder: capitalColor
  timeEncoder:
    layout: 2006-01-02 15:04:05.000
  durationEncoder: string
  callerEncoder: default
```

For more options see [config.example.yml][Config example].

By default, package search yaml-config in current dir or in `/etc/zap/config.yaml`.
Optionally there is ability to parametrize custom path to config using path-param in `zapper.New(nil, "/custom/path")`.

If zapper initialized with fallback logger, failures will be logged and Zap Logger will be initialized with default
configuration.

## Writer
Created for Zap Logger usage in [net/http][net/http] as ErrorLogger:
```go
package main

import (
	"log"
	"net/http"

	"github.com/nafigator/zapper"
	"github.com/nafigator/zapper/writer"
)

func main() {
	zl := zapper.Must(nil, nil)
	
	api := &http.Server{
		ErrorLog: log.New(writer.New(zl), "", 0),
	}

	if err := api.ListenAndServe(); err != nil {
		zl.Fatal("Http server failure: ", err)
	}
}
```

[Release img]: https://img.shields.io/badge/release-0.2.0-red.svg
[Release src]: https://github.com/nafigator/zap
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg
[Config example]: https://github.com/nafigator/zapper/blob/main/config.example.yml
[net/http]: https://pkg.go.dev/net/http
