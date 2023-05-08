# Uber Zap logger helpers
[![GitHub license][License img]][License src] [![GitHub release][Release img]][Release src] [![Scrutinizer master status][Scrutinizer master status image]][Scrutinizer master status src] [![Scrutinizer master code quality][Scrutinizer master quality image]][Scrutinizer master src] [![Conventional Commits][Conventional commits badge]][Conventional commits src]

## Conf
Features:
- easiest Zap Logger configuration using yaml-config.
- configuration TCP and UDP Zap Logger sinks out of box.
- minimal code for full-featured Zap Logger initialization. 

### Getting started
1. Install package:
    ```sh
    go install github.com/nafigator/zapper@latest
    ```
2. Use zapper for Zap Logger initialization:
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

By default, package search yaml-config in current dir or in `/etc/zap/config.yml`.
Optionally there is ability to parametrize custom path to config using path-param in `zapper.New("/custom/path", nil)`.

If zapper initialized with fallback logger, config read failures will be logged and Zap Logger will be initialized with
default configuration.

## Writer
Created for Zap Logger usage in [net/http][net/http] as ErrorLogger:
### Getting started
1. Install package
	```shell
	go install github.com/nafigator/zapper/writer@latest
	```
2. Use writer to initialize Zap Logger writer for http-server `ErrorLog`:
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

## Versioning
This software follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

[License img]: https://img.shields.io/badge/license-MIT-brightgreen.svg
[License src]: https://www.tldrlegal.com/license/mit-license
[Release img]: https://img.shields.io/badge/release-0.4.0-red.svg
[Release src]: https://github.com/nafigator/zap
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-blue.svg
[Config example]: https://github.com/nafigator/zapper/blob/main/config.example.yml
[net/http]: https://pkg.go.dev/net/http
[semver src]: http://semver.org
[Scrutinizer master quality image]: https://scrutinizer-ci.com/g/nafigator/zapper/badges/quality-score.png?b=master
[Scrutinizer master src]: https://scrutinizer-ci.com/g/nafigator/zapper/?branch=master
[Scrutinizer master status image]: https://scrutinizer-ci.com/g/nafigator/zapper/badges/build.png?b=master
[Scrutinizer master status src]: https://scrutinizer-ci.com/g/nafigator/zapper/?branch=master
