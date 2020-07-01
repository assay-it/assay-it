package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

type tOpts struct {
	secret *string
	source *string
	branch *string
	commit *string
	title  *string
	number *string
	api    *string
}

var opts *tOpts = &tOpts{}

func (opts *tOpts) parse() {
	opts.secret = flag.String("secret", "", "secret token for progammable access")
	opts.source = flag.String("source", "", "name of source code repository for assay suites (e.g. assay-it/example.assay.it).")
	opts.branch = flag.String("branch", "", "name of the branch with assay suites.")
	opts.commit = flag.String("commit", "", "long identity of commit.")
	opts.title = flag.String("title", "", "short description about the change (e.g. pull request title).")
	opts.number = flag.String("number", "", "pull request number.")
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

type Hook struct {
	ID          string       `json:"id"`
	PullRequest *PullRequest `json:"request,omitempty"`
}

func main() {
	log.SetFlags(0)
	flag.Usage = opts.usage
	opts.parse()

	if *opts.secret == "" {
		fmt.Fprintf(os.Stderr, "No security token.\n")
		os.Exit(1)
	}

	if *opts.source == "" {
		fmt.Fprintf(os.Stderr, "Undefined source code repository.\n")
		os.Exit(1)
	}

	if *opts.branch != "" && *opts.commit == "" {
		fmt.Fprintf(os.Stderr, "Undefined commit.\n")
		os.Exit(1)
	}

	if *opts.branch == "" && *opts.commit != "" {
		fmt.Fprintf(os.Stderr, "Undefined branch.\n")
		os.Exit(1)
	}

	//
	var token string

	category := "webhook"
	target := strings.Join(
		append([]string{"github"}, strings.Split(*opts.source, "/")...),
		":",
	)

	if *opts.branch != "" && *opts.commit != "" {
		category = "commit"
		target = fmt.Sprintf("%s:%s:%s", target, *opts.branch, *opts.commit)
	}

	hook := Hook{ID: target}
	if *opts.number != "" && *opts.title != "" {
		hook.PullRequest = &PullRequest{
			Number: *opts.number,
			Title:  *opts.title,
		}
	}

	err := gurl.Join(
		ioAccessToken(*opts.api, *opts.secret, &token),
		ioWebHook(*opts.api, category, hook, &token),
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
		ƒ.Code(200),
		ƒ.Recv(&rsp),
		ƒ.FMap(func() error {
			*token = "Bearer " + rsp.Token
			return nil
		}),
	)
}

func ioWebHook(api, cat string, req Hook, token *string) gurl.Arrow {
	var hook []byte

	return gurl.HTTP(
		ø.POST("https://%s/webhook/%s", api, cat),
		ø.Authorization().Val(token),
		ø.ContentJSON(),
		ø.Send(req),
		ƒ.Code(200, 401),
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
