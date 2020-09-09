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

	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

//
func (c *Client) WebHookSourceCode(sourcecode, target string) gurl.Arrow {
	var hook []byte

	return gurl.HTTP(
		ø.POST("https://%s/webhook/sourcecode", c.api),
		ø.Authorization().Val(&c.token),
		ø.ContentJSON(),
		ø.Send(SourceCodeID{
			ID:  sourcecode,
			URL: target,
		}),
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

//
func (c *Client) WebHook(req Hook) gurl.Arrow {
	var hook []byte

	return gurl.HTTP(
		ø.POST("https://%s/webhook/commit", c.api),
		ø.Authorization().Val(&c.token),
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
