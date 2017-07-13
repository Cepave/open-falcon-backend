package testing

import (
	"gopkg.in/check.v1"
	"net/url"
)

func ParseUrl(c *check.C, rawurl string) *url.URL {
	urlObject, err := url.Parse(rawurl)
	c.Assert(err, check.IsNil)

	return urlObject
}
func ParseRequestUri(c *check.C, rawurl string) *url.URL {
	urlObject, err := url.ParseRequestURI(rawurl)
	c.Assert(err, check.IsNil)

	return urlObject
}
func ParseQuery(c *check.C, query string) url.Values {
	values, err := url.ParseQuery(query)
	c.Assert(err, check.IsNil)

	return values
}
