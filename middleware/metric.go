package middleware

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
)

type metricCtxKey string

const fileServerHitsKey metricCtxKey = "fileServerHitsKey"

type MetricMiddleware struct {
	fileserverHits atomic.Int32
}

func NewMetricMiddleware() *MetricMiddleware {
	return &MetricMiddleware{}
}

func (m *MetricMiddleware) FileServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.fileserverHits.Add(1)
		ctx := context.WithValue(r.Context(), fileServerHitsKey, &m.fileserverHits)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetFileServerHitsFromContext(ctx context.Context) (*atomic.Int32, error) {
	fileserverHits, ok := ctx.Value(fileServerHitsKey).(*atomic.Int32)
	if !ok {
		return nil, errors.New("file server hits not found in context")
	}
	return fileserverHits, nil
}
