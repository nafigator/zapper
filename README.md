# Uber Zap logger helpers
[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

## Zapper
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

## Versioning
This software follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

[License img]: https://img.shields.io/badge/license-MIT-brightgreen.svg
[License src]: https://www.tldrlegal.com/license/mit-license
[Release img]: https://img.shields.io/badge/release-0.4.3-red.svg
[Release src]: https://github.com/nafigator/zapper
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-blue.svg
[Config example]: https://github.com/nafigator/zapper/blob/main/config.example.yml
[semver src]: http://semver.org
[Github main status src]: https://github.com/nafigator/zapper/tree/main
[Github main status badge]: https://github.com/nafigator/zapper/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/zapper
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/zapper
[Codecov report src]: https://app.codecov.io/gh/nafigator/zapper/tree/main
[Codecov report badge]: https://codecov.io/gh/nafigator/zapper/branch/main/graph/badge.svg
