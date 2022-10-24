package testhelper

import "net/http"

type AssertRequest func(r *http.Request)
type AssertingRoundTripper struct {
	Err      error
	Response *http.Response
	Assert   AssertRequest
}

func (r *AssertingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if r.Assert != nil {
		r.Assert(req)
	}
	return r.Response, r.Err
}

func NewTestClient(rt http.RoundTripper) *http.Client {
	return &http.Client{
		Transport: rt,
	}
}
