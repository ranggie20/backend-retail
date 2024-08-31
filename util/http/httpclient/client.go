package httpclient

import (
	"net"
	"net/http"
	"time"
)

type option struct {
	httpTransport http.RoundTripper
}

type Option func(*option)

// WithHTTPTransport set base transport before optionally traced. e.g., logged transport.
func WithHTTPTransport(t http.RoundTripper) Option { return func(o *option) { o.httpTransport = t } }

type Config struct {
	DialTimeout       time.Duration
	ConnectionTimeout time.Duration
	IdleTimeout       time.Duration
	MaxConn           uint
	MaxIdleConn       uint
}

// NewTransport create a new HTTP round tripper within the configured values.
func (c Config) NewTransport() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: c.DialTimeout,
		}).DialContext,
		MaxConnsPerHost:       int(c.MaxConn),
		MaxIdleConnsPerHost:   int(c.MaxIdleConn),
		IdleConnTimeout:       c.IdleTimeout,
		TLSHandshakeTimeout:   c.DialTimeout,
		ExpectContinueTimeout: c.DialTimeout,
	}
}

func New(cfg Config, opts ...Option) *http.Client {
	opt := &option{httpTransport: nil}
	for _, o := range opts {
		o(opt)
	}
	if opt.httpTransport == nil {
		opt.httpTransport = cfg.NewTransport()
	}
	return &http.Client{Transport: opt.httpTransport, Timeout: cfg.ConnectionTimeout}
}
