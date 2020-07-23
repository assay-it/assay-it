package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

type tOpts struct {
	secret *string
	head   *string
	base   *string
	title  *string
	number *string
	hub    *string
	api    *string
}

var opts *tOpts = &tOpts{}

func (opts *tOpts) parse() {
	opts.secret = flag.String("secret", "", "secret token for progammable access")
	opts.head = flag.String("head", "", "head of pull request such as owner/repo/branch/commit.")
	opts.base = flag.String("base", "", "base of pull request such as owner/repo/branch/commit.")
	opts.title = flag.String("title", "", "short description about the change (e.g. pull request title).")
	opts.number = flag.String("number", "", "pull request number.")
	opts.hub = flag.String("hub", "github", "hub service")
	opts.api = flag.String("api", "api.assay.it", "rest api endpoint.")
	flag.Parse()
}

func (opts *tOpts) usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "the application triggers a formal (objective) proofs of the quality using Behavior as a Code paradigm at https://assay.it\n")
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

type PullRequest struct {
	Number string `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
}

type Commit struct {
	ID string `json:"id"`
}

type Hook struct {
	PullRequest *PullRequest `json:"request,omitempty"`
	Base        Commit       `json:"base"`
	Head        Commit       `json:"head"`
}

func main() {
	log.SetFlags(0)
	flag.Usage = opts.usage
	opts.parse()

	if *opts.secret == "" {
		fmt.Fprintf(os.Stderr, "No security token.\n")
		os.Exit(1)
	}

	if *opts.base == "" || *opts.head == "" {
		fmt.Fprintf(os.Stderr, "No reference to changes.\n")
		os.Exit(1)
	}

	//
	var token string

	base := strings.Join(strings.Split(filepath.Join(*opts.hub, *opts.base), "/"), ":")
	head := strings.Join(strings.Split(filepath.Join(*opts.hub, *opts.head), "/"), ":")

	hook := Hook{
		Base: Commit{ID: base},
		Head: Commit{ID: head},
	}
	if *opts.number != "" && *opts.title != "" {
		hook.PullRequest = &PullRequest{
			Number: *opts.number,
			Title:  *opts.title,
		}
	}

	err := gurl.Join(
		ioAccessToken(*opts.api, *opts.secret, &token),
		ioWebHook(*opts.api, hook, &token),
	)(gurl.IO()).Fail

	if err != nil {
		fmt.Printf("Unable to trigger webhook: %v", err)
		os.Exit(1)
	}
}

func ioAccessToken(api, digest string, token *string) gurl.Arrow {
	req := struct {
		Type string `json:"grant_type"`
	}{Type: "client_credentials"}
	var rsp struct {
		Token string `json:"access_token"`
	}

	return gurl.HTTP(
		ø.POST("https://%s/auth/token", api),
		ø.Authorization().Is("Basic "+digest),
		ø.ContentForm(),
		ø.Send(req),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&rsp),
		ƒ.FMap(func() error {
			*token = "Bearer " + rsp.Token
			return nil
		}),
	)
}

func ioWebHook(api string, req Hook, token *string) gurl.Arrow {
	var hook []byte

	return gurl.HTTP(
		ø.POST("https://%s/webhook/commit", api),
		ø.Authorization().Val(token),
		ø.ContentJSON(),
		ø.Send(req),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Bytes(&hook),
		ƒ.FMap(func() error {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, hook, "", "  "); err != nil {
				return err
			}
			fmt.Println(pretty.String())
			return nil
		}),
	)
}

