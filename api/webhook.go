//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package api

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/assay-it/sdk-go/assay"
	ç "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

//
func (c *Client) WebHookSourceCode(sourcecode, target string) assay.Arrow {
	var hook []byte

	return http.Join(
		ø.POST("https://%s/webhook/sourcecode", c.api),
		ø.Authorization().Val(&c.token),
		ø.ContentJSON(),
		ø.Send(SourceCodeID{
			ID:  sourcecode,
			URL: target,
		}),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Bytes(&hook),
	).Then(
		ç.FMap(func() error {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, hook, "", "  "); err != nil {
				return err
			}
			fmt.Println(pretty.String())
			return nil
		}),
	)
}

//
func (c *Client) WebHook(req Hook) assay.Arrow {
	var hook []byte

	return http.Join(
		ø.POST("https://%s/webhook/commit", c.api),
		ø.Authorization().Val(&c.token),
		ø.ContentJSON(),
		ø.Send(req),
		ƒ.Code(http.StatusCodeOK),
		ƒ.Bytes(&hook),
	).Then(
		ç.FMap(func() error {
			var pretty bytes.Buffer
			if err := json.Indent(&pretty, hook, "", "  "); err != nil {
				return err
			}
			fmt.Println(pretty.String())
			return nil
		}),
	)
}
