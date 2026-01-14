package middleware

import (
	"net/http"
	"time"

	"github.com/maqsatto/Notes-API/internal/logger"
)

type statusWriter struct {
	status int
	bytes int
	http.ResponseWriter
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusWriter)  Write(b []byte) (int, error) {
	// if handler never called WriteHeader, status is 200 by default
	if w.status == 0 {
		w.status = http.StatusOK
	}

	n, err := w.ResponseWriter.Write(b)
	w.bytes += n
	return n, err
}

func Logger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			sw := *&statusWriter{ResponseWriter: w}
			next.ServeHTTP(&sw, r)
			duration := time.Since(start)
			log.Info(
				"request " +
					"method=" + r.Method +
					"path=" + r.URL.Path +
					"status=" + itoa(sw.status) +
					"bytes=" + itoa(sw.bytes) +
					"dur=" + duration.String(),
			)
		})
	}
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	// quick conversion
	buf := make([]byte, 0, 12)
	neg := v < 0
	if neg {
		v = -v
	}
	for v > 0 {
		buf = append(buf, byte('0'+v%10))
		v /= 10
	}
	if neg {
		buf = append(buf, '-')
	}
	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}
	return string(buf)
}
