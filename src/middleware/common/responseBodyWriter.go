package middleware

import (
	"bufio"
	"bytes"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	noWritten     = -1
	defaultStatus = 200
)

// https://github.com/jim3ma/gin-jsonp/blob/master/jsonp.go#L30
type responseBuffer struct {
	Response gin.ResponseWriter // the actual ResponseWriter to flush to
	status   int                // the HTTP response code from WriteHeader
	Body     *bytes.Buffer      // the response content body
	Flushed  bool
}

func NewResponseBuffer(w gin.ResponseWriter) *responseBuffer {
	return &responseBuffer{
		Response: w, status: defaultStatus, Body: &bytes.Buffer{},
	}
}

func (w *responseBuffer) Header() http.Header {
	return w.Response.Header() // use the actual response header
}

func (w *responseBuffer) Write(buf []byte) (int, error) {
	w.Body.Write(buf)
	return len(buf), nil
}

func (w *responseBuffer) WriteString(s string) (n int, err error) {
	//w.WriteHeaderNow()
	//n, err = io.WriteString(w.ResponseWriter, s)
	//w.size += n
	n, err = w.Write([]byte(s))
	return
}

func (w *responseBuffer) Written() bool {
	return w.Body.Len() != noWritten
}

func (w *responseBuffer) WriteHeader(status int) {
	w.status = status
}

func (w *responseBuffer) WriteHeaderNow() {
	//if !w.Written() {
	//	w.size = 0
	//	w.ResponseWriter.WriteHeader(w.status)
	//}
}

func (w *responseBuffer) Status() int {
	return w.status
}

func (w *responseBuffer) Size() int {
	return w.Body.Len()
}

func (w *responseBuffer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	//if w.size < 0 {
	//	w.size = 0
	//}
	return w.Response.(http.Hijacker).Hijack()
}

func (w *responseBuffer) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

// Fake Flush
// TBD
func (w *responseBuffer) Flush() {
	w.realFlush()
}

func (w *responseBuffer) realFlush() {
	if w.Flushed {
		return
	}
	w.Response.WriteHeader(w.status)
	if w.Body.Len() > 0 {
		_, err := w.Response.Write(w.Body.Bytes())
		if err != nil {
			panic(err)
		}
		w.Body.Reset()
	}
	w.Flushed = true
}

func (w *responseBuffer) Pusher() http.Pusher {
	return w.Response.Pusher()
}
