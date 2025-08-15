package ssh

import "time"

type Option interface {
	Apply(*Client)
}

type OptionFunc func(*Client)

func (f OptionFunc) Apply(c *Client) {
	f(c)
}

func WithRetries(n int) Option {
	return OptionFunc(func(c *Client) {
		c.retries = n
	})
}

func WithTimeout(d time.Duration) Option {
	return OptionFunc(func(c *Client) {
		c.timeout = d
	})
}
