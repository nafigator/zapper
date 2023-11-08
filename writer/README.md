# Writer
[![GitHub release][Release img]][Release src] [![Github main status][Github main status badge]][Github main status src] [![Go Report Card][Go Report Card badge]][Go Report Card src] [![Coverage report][Codecov report badge]][Codecov report src]

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
		"github.com/nafigator/zapper/conf"
		"github.com/nafigator/zapper/writer"
	)

	func main() {
		zl := zapper.Must(conf.Must())

		api := &http.Server{
			ErrorLog: log.New(writer.New(zl), "", 0),
		}

		if err := api.ListenAndServe(); err != nil {
			zl.Fatal("Http server failure: ", err)
		}
	}
	```

### Versioning
This software follows *"Semantic Versioning"* specifications. The signature of exported package functions is used
as a public API. Read more on [SemVer.org][semver src].

[License img]: https://img.shields.io/badge/license-MIT-brightgreen.svg
[License src]: https://www.tldrlegal.com/license/mit-license
[Release img]: https://img.shields.io/badge/release-0.1.1-red.svg
[Release src]: https://github.com/nafigator/zapper/tree/main/writer
[Conventional commits src]: https://conventionalcommits.org
[Conventional commits badge]: https://img.shields.io/badge/Conventional%20Commits-1.0.0-blue.svg
[net/http]: https://pkg.go.dev/net/http
[semver src]: http://semver.org
[Github main status src]: https://github.com/nafigator/zapper/tree/main
[Github main status badge]: https://github.com/nafigator/zapper/actions/workflows/go.yml/badge.svg?branch=main
[Go Report Card src]: https://goreportcard.com/report/github.com/nafigator/zapper/writer
[Go Report Card badge]: https://goreportcard.com/badge/github.com/nafigator/zapper/writer
[Codecov report src]: https://app.codecov.io/gh/nafigator/zapper/tree/main/writer
[Codecov report badge]: https://codecov.io/gh/nafigator/zapper/writer/branch/coverage/graph/badge.svg
