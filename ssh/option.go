package ssh

import "time"

// Option configures an SSH Client.
type Option interface {
	Apply(*Client)
}

// OptionFunc is an adapter to allow ordinary functions to be used as Options.
type OptionFunc func(*Client)

// Apply calls the function with the given Client.
func (f OptionFunc) Apply(c *Client) {
	f(c)
}

// WithRetries sets the number of SSH connection retry attempts.
func WithRetries(n int) Option {
	return OptionFunc(func(c *Client) {
		c.retries = n
	})
}

// WithTimeout sets the SSH connection timeout.
func WithTimeout(d time.Duration) Option {
	return OptionFunc(func(c *Client) {
		c.timeout = d
	})
}
