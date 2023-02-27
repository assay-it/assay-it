//
// Copyright (C) 2020 - 2023 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package httpbin

import (
	"github.com/fogfish/gurl/v2/http"
	ƒ "github.com/fogfish/gurl/v2/http/recv"
	ø "github.com/fogfish/gurl/v2/http/send"
)

func TestHttpBinGet() http.Arrow {
	return http.GET(
		ø.URI("http://httpbin.org/get"),
		ø.UserAgent.Set("curl/7.64.1"),
		ø.Accept.Set("*/*"),

		ƒ.Status.OK,
		ƒ.ContentType.ApplicationJSON,
		ƒ.Header("Access-Control-Allow-Origin", "*"),
		ƒ.Header("Access-Control-Allow-Credentials", "true"),
		ƒ.Match(`
			{
				"args": {},
				"headers": {
					"Accept": "*/*",
					"Host": "httpbin.org",
					"User-Agent": "curl/7.64.1"
				},
				"origin": "_",
				"url": "http://httpbin.org/get"
			}
		`),
	)
}
