<p align="center">
  <img src="./doc/assay-it.svg" height="120" />
  <h3 align="center">assay-it</h3>
  <p align="center"><strong>Test Microservice in Production. Confirm Quality & Eliminate Risk.</strong></p>

  <p align="center">
    <!-- Discussion -->
    <a href="https://github.com/assay-it/assay-it/discussions">
      <img alt="GitHub Discussions" src="https://img.shields.io/github/discussions/assay-it/assay-it?logo=github">
    </a>
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

[User Guide](https://assay.it/doc/) |
[Community Support](https://github.com/assay-it/assay-it/discussions) |
[Golang Suite example](./examples/golang-httpbin/request.go) |
[Markdown Suite example](./examples/katt-httpbin/request.md)

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

## How To Contribute

The command line utility is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request


The build and testing process requires [Go](https://golang.org) version 1.20 or later.

**Build** and **run** in your development console.

```bash
git clone https://github.com/assay-it/assay-it
cd assay-it
go test
```

### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/assay-it/assay-it/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/assay-it/assay-it.svg?style=for-the-badge)](LICENSE)
