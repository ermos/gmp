package gmp

import (
	"net/url"
)

func (c *Client) buildUrl(base string, query map[string]string) string {
	v := url.Values{}

	for key, value := range query {
		v.Set(key, value)
	}
	v.Set("key", c.apiKey)

	return base + "?" + v.Encode()
}