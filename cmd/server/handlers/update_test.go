package handlers

import (
	metric "alerting/internal/metrics"
	mstorage "alerting/internal/mstorage"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricHandler_UpdateRequest(t *testing.T) {

	var MemStorage mstorage.MetricStorage = mstorage.InMemoryStorage{
		Storage: map[string]metric.AbstractMetric{}}

	mHandler := MetricHandler{
		Storage: &MemStorage,
	}

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", mHandler.UpdateRequest)
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Get("/", mHandler.ValueAllMetrics)
		r.Route("/value/{metricType}/{metricName}", func(r chi.Router) {
			r.Get("/", mHandler.ValueMetricByName)
		})
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testCases := []struct {
		method       string
		requestURL   string
		expectedCode int
	}{
		{method: http.MethodPost, requestURL: "/update/gauge/test/1", expectedCode: http.StatusOK},
		{method: http.MethodPost, requestURL: "/update/test/test/1", expectedCode: http.StatusBadRequest},
		{method: http.MethodPost, requestURL: "/update/gauge//1", expectedCode: http.StatusNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = ts.URL + tc.requestURL

			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")

		})
	}
}
