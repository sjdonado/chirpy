package middleware

import (
	"net/http"
	"sync/atomic"
)

type MetricMiddleware struct {
	fileserverHits atomic.Int32
}

func NewMetricMiddleware() *MetricMiddleware {
	return &MetricMiddleware{}
}

func (m *MetricMiddleware) FileServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (m *MetricMiddleware) GetFileServerHits() int32 {
	return m.fileserverHits.Load()
}

func (m *MetricMiddleware) ResetFileServerHits() {
	m.fileserverHits.Store(0)
}
