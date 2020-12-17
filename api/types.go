//
// Copyright (C) 2020 assay.it
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/assay.it/assay
//

package api

//
type Client struct {
	api   string
	token string
}

//
type PullRequest struct {
	Number string `json:"number,omitempty"`
	Title  string `json:"title,omitempty"`
}

//
type SourceCodeID struct {
	ID          string       `json:"id"`
	PullRequest *PullRequest `json:"request,omitempty"`
	URL         string       `json:"endpoint,omitempty"`
}

//
type Commit struct {
	ID string `json:"id"`
}

//
type Hook struct {
	PullRequest *PullRequest `json:"request,omitempty"`
	Base        Commit       `json:"base"`
	Head        Commit       `json:"head"`
	URL         string       `json:"endpoint,omitempty"`
}

//
func New(api string) *Client {
	return &Client{api: api}
}
