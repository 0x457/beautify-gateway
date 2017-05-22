package proxy

import (
	"net/http"
	"sync"
)

// ResponseWriter represents HTTP response.
//
// It is forbidden copying Response instances. Create new instances
// and use CopyTo instead.
//
// Response instance MUST NOT be used from concurrently running goroutines.
type ResponseWriter interface {
	http.ResponseWriter
	releaseMe()
}

var rpool = sync.Pool{New: func() interface{} { return &responseWriter{statusCode: StatusOK} }}

func acquireResponseWriter(underline http.ResponseWriter) *responseWriter {
	w := rpool.Get().(*responseWriter)
	w.ResponseWriter = underline
	return w
}

func releaseResponseWriter(w *responseWriter) {
	w.statusCodeSent = false
	w.beforeFlush = nil
	w.statusCode = StatusOK
	rpool.Put(w)
}

// responseWriter is the basic response writer,
// it writes directly to the underline http.ResponseWriter
type responseWriter struct {
	http.ResponseWriter
	statusCode     int  // the saved status code which will be used from the cache service
	statusCodeSent bool // reply header has been (logically) written
	// yes only one callback, we need simplicity here because on EmitError the beforeFlush events should NOT be cleared
	// but the response is cleared.
	// Sometimes is useful to keep the event,
	// so we keep one func only and let the user decide when he/she wants to override it with an empty func before the EmitError (context's behavior)
	beforeFlush func()
}

// var _ ResponseWriter = &responseWriter{}

// StatusCode returns the status code header value
func (w *responseWriter) StatusCode() int {
	return w.statusCode
}

// Write writes to the client
// If WriteHeader has not yet been called, Write calls
// WriteHeader(http.StatusOK) before writing the data. If the Header
// does not contain a Content-Type line, Write adds a Content-Type set
// to the result of passing the initial 512 bytes of written data to
// DetectContentType.
//
// Depending on the HTTP protocol version and the client, calling
// Write or WriteHeader may prevent future reads on the
// Request.Body. For HTTP/1.x requests, handlers should read any
// needed request body data before writing the response. Once the
// headers have been flushed (due to either an explicit Flusher.Flush
// call or writing enough data to trigger a flush), the request body
// may be unavailable. For HTTP/2 requests, the Go HTTP server permits
// handlers to continue to read the request body while concurrently
// writing the response. However, such behavior may not be supported
// by all HTTP/2 clients. Handlers should read before writing if
// possible to maximize compatibility.
func (w *responseWriter) Write(contents []byte) (int, error) {
	w.tryWriteHeader()
	return w.ResponseWriter.Write(contents)
}

func (w *responseWriter) tryWriteHeader() {
	if !w.statusCodeSent { // by write
		w.statusCodeSent = true
		w.ResponseWriter.WriteHeader(w.statusCode)
	}
}
func (w *responseWriter) releaseMe() {
	releaseResponseWriter(w)
}
