<p align="center">
  <img src="./doc/assay-it.svg" height="120" />
  <h3 align="center">assay-it</h3>
  <p align="center"><strong>Test Microservice in Production. Confirm Quality & Eliminate Risk.</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/assay-it/assay-it/releases">
      <img src="https://img.shields.io/github/v/tag/assay-it/assay-it?label=version" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/assay-it/assay-it/actions/">
      <img src="https://github.com/assay-it/assay-it/workflows/test/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/assay-it/assay-it">
      <img src="https://img.shields.io/github/last-commit/assay-it/assay-it.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/assay-it/assay-it?branch=main">
      <img src="https://coveralls.io/repos/github/assay-it/assay-it/badge.svg?branch=main" />
    </a>
  </p>
</p>

--- 

Construct automated quality check pipelines for your applications, microservices and other endpoints across deployments and environments.

[![asciicast](https://asciinema.org/a/564197.svg)](https://asciinema.org/a/564197)


## Quick Example

First install the assay-it command line utility

```bash 
brew tap assay-it/homebrew-tap
brew install -q assay-it
```

Then implement test scenario using [type safe, pure functional Golang combinators](https://github.com/fogfish/gurl) defined by ᵍ🆄🆁🅻 library. 

```go
// httpbin.go file
package httpbin

import (
  "github.com/fogfish/gurl/v2/http"
  ƒ "github.com/fogfish/gurl/v2/http/recv"
  ø "github.com/fogfish/gurl/v2/http/send"
)

func TestHttpBinGet() http.Arrow {
  return http.GET(
    ø.URI("http://httpbin.org/get"),

    ƒ.Status.OK,
    ƒ.ContentType.ApplicationJSON,
  )
}
```

Now start testing

```bash
assay-it test httpbin.go
```


## License

[![See LICENSE](https://img.shields.io/github/license/assay-it/assay.svg?style=for-the-badge)](LICENSE)
