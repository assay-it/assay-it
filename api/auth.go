//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package api

import (
	"github.com/fogfish/gurl"
	ƒ "github.com/fogfish/gurl/http/recv"
	ø "github.com/fogfish/gurl/http/send"
)

type tOAuthFlow struct {
	Type string `json:"grant_type"`
}

type tAccessToken struct {
	Token string `json:"access_token"`
}

//
func (c *Client) SignIn(digest string) gurl.Arrow {
	var token tAccessToken

	return gurl.HTTP(
		ø.POST("https://%s/auth/token", c.api),
		ø.Authorization().Is("Basic "+digest),
		ø.ContentForm(),
		ø.Send(tOAuthFlow{Type: "client_credentials"}),
		ƒ.Code(gurl.StatusCodeOK),
		ƒ.Recv(&token),
		ƒ.FMap(func() error {
			c.token = "Bearer " + token.Token
			return nil
		}),
	)
}
