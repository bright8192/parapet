package redirect

import (
	"net/http"
	"strings"

	"github.com/moonrhythm/parapet/pkg/header"
)

// NonWWW creates new non www redirector
func NonWWW() *NonWWWRedirector {
	return new(NonWWWRedirector)
}

// NonWWWRedirector redirects to non-www
type NonWWWRedirector struct {
	StatusCode int
}

// ServeHandler implements middleware interface
func (m NonWWWRedirector) ServeHandler(h http.Handler) http.Handler {
	if m.StatusCode <= 0 {
		m.StatusCode = http.StatusMovedPermanently
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.TrimPrefix(r.Host, "www.")
		if len(host) < len(r.Host) {
			proto := header.Get(r.Header, header.XForwardedProto)
			http.Redirect(w, r, proto+"://"+host+r.RequestURI, m.StatusCode)
			return
		}
		h.ServeHTTP(w, r)
	})
}
