package mixed

import "net/http"

// tempResponseWriter implements the http.ResponseWriter interface
// It used when you need to temporary hold the content and use it later
type tempResponseWriter struct {
	header     http.Header
	StatusCode int
	Content    []byte
}

func (rw *tempResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *tempResponseWriter) Write(content []byte) (int, error) {
	rw.Content = append(rw.Content, content...)
	return len(content), nil
}

func (rw *tempResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
}
