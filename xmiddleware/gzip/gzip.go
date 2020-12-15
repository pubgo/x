package gzip

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/pubgo/x/xconfig"
	"github.com/pubgo/x/xdi"
	"github.com/pubgo/x/xerror"
	"io/ioutil"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	BestCompression    = gzip.BestCompression
	BestSpeed          = gzip.BestSpeed
	DefaultCompression = gzip.DefaultCompression
	NoCompression      = gzip.NoCompression
)

func SetUp() gin.HandlerFunc {
	var level int
	xerror.Panic(xdi.Invoke(func(config *xconfig.Config) {
		level = config.Web.Gzip.Level
	}))

	var gzPool sync.Pool
	gzPool.New = func() interface{} {
		gz, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return gz
	}
	return func(c *gin.Context) {
		if !shouldCompress(c.Request) {
			return
		}

		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(c.Writer)

		c.Header("Content-Encoding", "gzip")
		// canGzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		c.Header("Vary", "Accept-Encoding")
		c.Writer = &gzipWriter{c.Writer, gz}
		defer func() {
			gz.Close()
			c.Header("Content-Length", fmt.Sprint(c.Writer.Size()))
		}()
		c.Next()
	}
}

type gzipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

func (g *gzipWriter) WriteString(s string) (int, error) {
	return g.writer.Write([]byte(s))
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.writer.Write(data)
}

// Fix: https://github.com/mholt/caddy/issues/38
func (g *gzipWriter) WriteHeader(code int) {
	g.Header().Del("Content-Length")
	g.ResponseWriter.WriteHeader(code)
}

// Provided in order to implement the http.Hijacker interface.
func (g *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker := g.ResponseWriter.(http.Hijacker)
	return hijacker.Hijack()
}

// Call the parent CloseNotify.
// Provided in order to implement the http.CloseNotifier interface.
func (g *gzipWriter) CloseNotify() <-chan bool {
	notifier := g.ResponseWriter.(http.CloseNotifier)
	return notifier.CloseNotify()
}

func shouldCompress(req *http.Request) bool {
	if !strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") ||
		strings.Contains(req.Header.Get("Connection"), "Upgrade") ||
		strings.Contains(req.Header.Get("Content-Type"), "text/event-stream") {
		return false
	}

	extension := filepath.Ext(req.URL.Path)
	if len(extension) < 4 { // fast path
		return true
	}

	switch extension {
	case ".png", ".gif", ".jpeg", ".jpg":
		return false
	default:
		return true
	}
}
