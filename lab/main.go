// Package main provides a simple web server that serves files from the dist
// directory, including the wasm payload.
package main

import (
	"compress/gzip"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var gz = flag.Bool("gzip", false, "enable automatic gzip compression")

func wasmContentTypeSetter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".wasm") {
			w.Header().Set("content-type", "application/wasm")
		}
		h.ServeHTTP(w, r)
	})
}

func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// Client does not accept gzip encoding. Serve the request as-is.
			h.ServeHTTP(w, r)
			return
		}
		// Client accepts gzip encoding. Compress the response using gzip.
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzResponseWriter := gzipResponseWriter{ResponseWriter: w, Writer: gz}
		h.ServeHTTP(&gzResponseWriter, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// loggingResponseWriter wraps http.ResponseWriter to capture the response size.
type loggingResponseWriter struct {
	http.ResponseWriter
	bytesWritten int
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += n
	return n, err
}

// loggingMiddleware logs the HTTP requests with the specified details.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get client IP address, considering proxy headers if present.
		ip := r.RemoteAddr
		if ipProxy := r.Header.Get("X-Forwarded-For"); ipProxy != "" {
			ips := strings.Split(ipProxy, ",")
			ip = strings.TrimSpace(ips[0])
		} else {
			// Remove port from r.RemoteAddr.
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				ip = host
			}
		}

		// Wrap the ResponseWriter to capture the response size.
		lrw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		sizeMB := float64(lrw.bytesWritten) / 1024.0 / 1024.0

		// Log the request details.
		log.Printf("%s %s %s %.2fMB %dms",
			ip,
			r.Method,
			r.URL.Path,
			sizeMB,
			duration.Milliseconds())
	})
}

func main() {
	flag.Parse()
	h := http.FileServer(http.Dir("./dist"))
	h = loggingMiddleware(h)
	h = wasmContentTypeSetter(h)
	if *gz {
		h = gzipHandler(h)
	}
	http.Handle("/", h) // Serve the files in the dist directory with logging.

	port := os.Getenv("PORT") // Set by heroku for deployment.
	if port == "" {
		port = "8060"
	}
	log.Print("Serving on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
