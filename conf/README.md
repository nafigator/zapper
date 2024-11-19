# Conf
[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

Created for simple Zap Logger configuration build based on yml-config.
### Getting started
1. Install package
	```shell
	go install github.com/nafigator/zapper/conf@latest
	```
2. Use conf to initialize Zap config:
	```go
	package main

	import (
		"github.com/nafigator/zapper"
		"github.com/nafigator/zapper/conf"
	)

	func main() {
 		// Initializes config based on yml-files /etc/zap/config.yml or ./config.yml
		log := zapper.Must(conf.Must())

		log.Info("Zap Logger initialized")
	}
	```

### Versioning
This software follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

[Release img]: https://img.shields.io/github/v/tag/nafigator/zapper?logo=github&color=teal&filter=conf*
[Release src]: https://github.com/nafigator/zapper/tree/main/conf
[semver src]: http://semver.org
[Github main status src]: https://github.com/nafigator/zapper/tree/main
[Github main status badge]: https://github.com/nafigator/zapper/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/zapper/conf
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/zapper/conf
[Codecov report src]: https://app.codecov.io/gh/nafigator/zapper/tree/main/conf
[Codecov report badge]: https://codecov.io/gh/nafigator/zapper/conf/branch/main/graph/badge.svg
