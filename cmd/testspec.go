package cmd

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	rootCmd.AddCommand(testSpecCmd)
	testSpecCmd.Flags().StringVarP(&testSpecFormat, "format", "f", "", "generate test suite in the format, supported formats go, md)")
}

var (
	testSpecFormat string
)

var testSpecCmd = &cobra.Command{
	Use:   "testspec",
	Short: "generate test specification",
	Long: `
Testing suites are type safe, pure functional Golang source code files. Each
file contains one or multiple test functions named TestXxx (where Xxx does not
start with a lower case letter) and should have the signature,
	
	func TestXxx() http.Arrow { ... } 

These functions declares cause-and-effect for protocols operations:
* "Given" specifies the communication context and the known state of
the expected behavior;
* "When" executes key actions about the interaction against
specified deployment;
* "Then" observes output, validates its correctness and outputs results.

	package suite

	import (
		"github.com/fogfish/gurl/v2/http"
		ƒ "github.com/fogfish/gurl/v2/http/recv"
		ø "github.com/fogfish/gurl/v2/http/send"
	)

	func TestXxx() http.Arrow {
		return http.GET(
			ø.URI("http://example.com"),
			ƒ.Status.OK,
		)
	}

See the documentation on https://assay.it for more information.
	`,
	Example: `
	assay-it testspec
	assay-it testspec -f md
	assay-it testspec -f go
	assay-it testspec -f go \
	  https://assay.it/doc/ \
	  https://assay.it/doc/introduction
	`,
	SilenceUsage: true,
	RunE:         testSpec,
}

func testSpec(cmd *cobra.Command, args []string) error {
	switch testSpecFormat {
	case "md":
		if len(args) == 0 {
			stdout.Write([]byte(testSpecMd()))
		}

		spec, err := testSpecMdFromUrls(args)
		if err != nil {
			return err
		}
		stdout.Write([]byte(spec))

	default:
		if len(args) == 0 {
			stdout.Write([]byte(testSpecGo()))
		}

		spec, err := testSpecGoFromUrls(args)
		if err != nil {
			return err
		}
		stdout.Write([]byte(spec))
	}
	return nil
}

func testSpecConfig() string {
	return `{
  "suites": [
    "suites/suite.go"
  ]
}`
}

func testSpecGo() string {
	return `
package suites

import (
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

func TestHttpBinGet() http.Arrow {
	return http.GET(
		ø.URI("http://httpbin.org/get"),
		ø.UserAgent.Set("gurl/v2"),
	
		ƒ.Status.OK,
		ƒ.ContentType.ApplicationJSON,
		ƒ.Match(` + "`" + `
			{
				"headers": {
					"Host": "httpbin.org",
					"User-Agent": "gurl/v2"
				},
				"origin": "_",
				"url": "http://httpbin.org/get"
			}
		` + "`" + `),
	)
}	
`
}

func testSpecMd() string {
	return `
## Test HttpBin Get

` + "```" + `
GET http://httpbin.org/get
> User-Agent: gurl/v2
< 200 OK
< Content-Type: application/json
{
  "headers": {
    "Host": "httpbin.org", 
    "User-Agent": "gurl/v2"
  }, 
  "origin": "_", 
  "url": "http://httpbin.org/get"
}
` + "```\n"
}

func testSpecGoFromUrls(urls []string) (string, error) {
	var s strings.Builder
	s.WriteString(`
package suites

import (
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)		

`)

	for _, raw := range urls {
		name, err := urlToGoName(raw)
		if err != nil {
			return "", err
		}

		s.WriteString(`
func Test` + name + `() http.Arrow {
	return http.GET(
		ø.URI("` + raw + `"),
		ø.UserAgent.Set("gurl/v2"),
		ƒ.Status.OK,
	)
}

`)
	}

	return s.String(), nil
}

func urlToGoName(raw string) (string, error) {
	uri, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	path := strings.ReplaceAll(uri.Path, string(filepath.Separator), " ")
	name := cases.Title(language.English).String(path)

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}

	return reg.ReplaceAllString(name, ""), nil
}

func testSpecMdFromUrls(urls []string) (string, error) {
	var s strings.Builder

	for _, raw := range urls {
		name, err := urlToMdName(raw)
		if err != nil {
			return "", err
		}

		s.WriteString(`
## Test ` + name + `

` + "```" + `
GET ` + raw + `
> User-Agent: gurl/v2
< 200 OK
` + "```" + `

`)
	}

	return s.String(), nil

}

func urlToMdName(raw string) (string, error) {
	uri, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	path := strings.ReplaceAll(uri.Path, string(filepath.Separator), " ")
	path = strings.ReplaceAll(path, "-", " ")
	name := cases.Title(language.English).String(path)

	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+")
	if err != nil {
		return "", err
	}

	return reg.ReplaceAllString(name, ""), nil
}
