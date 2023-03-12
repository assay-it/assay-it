package cmd

import "github.com/spf13/cobra"

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
	`,
	SilenceUsage: true,
	RunE:         testSpec,
}

func testSpec(cmd *cobra.Command, args []string) error {
	switch testSpecFormat {
	case "md":
		stdout.Write([]byte(testSpecMd()))
	default:
		stdout.Write([]byte(testSpecGo()))
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
		ø.UserAgent.Set("curl/7.64.1"),
	
		ƒ.Status.OK,
		ƒ.ContentType.ApplicationJSON,
		ƒ.Match(` + "`" + `
			{
				"headers": {
					"Host": "httpbin.org",
					"User-Agent": "curl/7.64.1"
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
> User-Agent: curl/7.64.1
< 200 OK
< Content-Type: application/json
{
  "headers": {
    "Host": "httpbin.org", 
    "User-Agent": "curl/7.64.1"
  }, 
  "origin": "_", 
  "url": "http://httpbin.org/get"
}
` + "```\n"
}
