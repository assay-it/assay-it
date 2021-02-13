//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package api

import (
	"github.com/assay-it/sdk-go/assay"
	ç "github.com/assay-it/sdk-go/cats"
	"github.com/assay-it/sdk-go/http"
	ƒ "github.com/assay-it/sdk-go/http/recv"
	ø "github.com/assay-it/sdk-go/http/send"
)

type tOAuthFlow struct {
	Type string `json:"grant_type"`
}

type tAccessToken struct {
	Token string `json:"access_token"`
}

//
func (c *Client) SignIn(digest string) assay.Arrow {
	var token tAccessToken

	return http.Join(
		ø.POST("https://%s/auth/token", c.api),
		ø.Authorization().Is("Basic "+digest),
		ø.ContentForm(),
		ø.Send(tOAuthFlow{Type: "client_credentials"}),
		ƒ.Code(http.StatusOK),
		ƒ.Recv(&token),
	).Then(
		ç.FMap(func() error {
			c.token = "Bearer " + token.Token
			return nil
		}),
	)
}
