package block

import (
	"net/http"

	"github.com/moonrhythm/parapet"
)

// New creates news block
func New(match func(r *http.Request) bool) *Block {
	return &Block{
		Match: match,
	}
}

// Block is middleware block
type Block struct {
	Match func(r *http.Request) bool
	ms    parapet.Middlewares
}

// Use uses middleware
func (b *Block) Use(m parapet.Middleware) {
	b.ms.Use(m)
}

func (b *Block) UseFunc(m parapet.MiddlewareFunc) {
	b.Use(m)
}

// ServeHandler implements middleware interface
func (b *Block) ServeHandler(h http.Handler) http.Handler {
	next := b.ms.ServeHandler(http.NotFoundHandler())

	if b.Match == nil {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if b.Match(r) {
			next.ServeHTTP(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}
