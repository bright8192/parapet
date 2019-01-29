package transport

import (
	"crypto/tls"
	"net"
	"net/http"
	"sync"
	"time"
)

// HTTPS transport
type HTTPS struct {
	once sync.Once
	h    *http.Transport

	DialTimeout           time.Duration
	Timeout               time.Duration
	TCPKeepAlive          time.Duration
	DisableKeepAlives     bool
	MaxIdleConns          int
	IdleConnTimeout       time.Duration
	ResponseHeaderTimeout time.Duration
	VerifyCA              bool
}

// RoundTrip implement http.RoundTripper
func (t *HTTPS) RoundTrip(r *http.Request) (*http.Response, error) {
	t.once.Do(func() {
		if t.TCPKeepAlive == 0 {
			t.TCPKeepAlive = 10 * time.Minute
		}
		if t.MaxIdleConns == 0 {
			t.MaxIdleConns = 100
		}
		if t.IdleConnTimeout == 0 {
			t.IdleConnTimeout = 10 * time.Minute
		}

		t.h = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   t.DialTimeout,
				KeepAlive: t.TCPKeepAlive,
				DualStack: true,
			}).DialContext,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !t.VerifyCA,
			},
			DisableKeepAlives:     t.DisableKeepAlives,
			MaxIdleConnsPerHost:   t.MaxIdleConns,
			IdleConnTimeout:       t.IdleConnTimeout,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableCompression:    true,
			ResponseHeaderTimeout: t.ResponseHeaderTimeout,
		}
	})

	r.URL.Scheme = "https"
	return t.h.RoundTrip(r)
}