package headers

import "net/http"

// InterceptRequest creates new request interceptor
func InterceptRequest(f func(http.Header)) *RequestInterceptor {
	return &RequestInterceptor{Intercept: f}
}

// RequestInterceptor intercepts request's headers
type RequestInterceptor struct {
	Intercept func(http.Header)
}

// ServeHandler implements middleware interface
func (m RequestInterceptor) ServeHandler(h http.Handler) http.Handler {
	if m.Intercept == nil {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Intercept(r.Header)
		h.ServeHTTP(w, r)
	})
}

// ResponseInterceptFunc is the function for response's interceptor
type ResponseInterceptFunc func(w ResponseHeaderWriter)

// ResponseHeaderWriter type
type ResponseHeaderWriter interface {
	StatusCode() int
	Header() http.Header
	WriteHeader(statusCode int)
}

// InterceptResponse creates new response interceptor
func InterceptResponse(f ResponseInterceptFunc) *ResponseInterceptor {
	return &ResponseInterceptor{Intercept: f}
}

// ResponseInterceptor intercepts response's headers
type ResponseInterceptor struct {
	Intercept ResponseInterceptFunc
}

// ServeHandler implements middleware interface
func (m ResponseInterceptor) ServeHandler(h http.Handler) http.Handler {
	if m.Intercept == nil {
		return h
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw := interceptRW{
			ResponseWriter: w,
			f:              m.Intercept,
			status:         http.StatusOK,
		}
		defer nw.intercept()

		h.ServeHTTP(&nw, r)
	})
}

type interceptRW struct {
	http.ResponseWriter
	wroteHeader bool
	intercepted bool
	status      int
	f           ResponseInterceptFunc
}

func (w *interceptRW) intercept() {
	if w.intercepted {
		return
	}
	w.intercepted = true
	w.f(w)
}

func (w *interceptRW) WriteHeader(statusCode int) {
	if !w.intercepted {
		w.status = statusCode
		w.intercept()
	}
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *interceptRW) Write(p []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.ResponseWriter.Write(p)
}

// StatusCode returns status code
func (w *interceptRW) StatusCode() int {
	return w.status
}
