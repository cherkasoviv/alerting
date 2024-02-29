package handlers

import (
	"alerting/internal/config"
	"alerting/internal/metrics"
	"alerting/internal/mstorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func Test_valueHandler_GetJSON(t *testing.T) {
	cfg := config.ServerConfig{
		Host:            "",
		StoreInterval:   0,
		FileStoragePath: "",
		NeedToRestore:   false,
	}
	storage := mstorage.Initialize(&cfg)
	valueHandler := NewValueHandler(storage)
	gm := metrics.GaugeMetric{
		Metric: struct {
			Name  string
			Mtype metrics.MetricType
		}{Name: "testGauge", Mtype: metrics.Gauge},
		Value: 1,
	}
	cm := metrics.CounterMetric{
		Metric: struct {
			Name  string
			Mtype metrics.MetricType
		}{Name: "testCounter", Mtype: metrics.Counter},
		Value: 1,
	}
	storage.CreateOrUpdateMetric(&gm)
	storage.CreateOrUpdateMetric(&cm)

	r := chi.NewRouter()
	r.Route("/value", func(r chi.Router) {

		r.Post("/", valueHandler.GetJSON())

	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testCases := []struct {
		name                 string
		method               string
		requestURL           string
		expectedCode         int
		requestJSON          string
		expectedResponseJSON string
	}{
		{name: "success gauge", method: http.MethodPost, requestURL: "/value", expectedCode: http.StatusOK,
			requestJSON:          "{\"id\":\"testGauge\" , \"type\": \"gauge\"}",
			expectedResponseJSON: "{\"id\":\"testGauge\",\"type\":\"gauge\",\"value\":1}\n"},
		{name: "success counter", method: http.MethodPost, requestURL: "/value", expectedCode: http.StatusOK,
			requestJSON:          "{\"id\":\"testCounter\" , \"type\": \"counter\"}",
			expectedResponseJSON: "{\"id\":\"testCounter\",\"type\":\"counter\",\"delta\":1}\n"},
		{name: "error json", method: http.MethodPost, requestURL: "/value", expectedCode: http.StatusBadRequest,
			requestJSON:          "{\"name\":\"test\" ,\"type\":\"gauge\"}",
			expectedResponseJSON: ""},
		{name: "no metric in storage", method: http.MethodPost, requestURL: "/value", expectedCode: http.StatusNotFound,
			requestJSON:          "{\"id\":\"test1\" ,\"type\":\"gauge\"}",
			expectedResponseJSON: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := resty.New().R()
			req.Method = tc.method
			req.URL = ts.URL + tc.requestURL
			req.SetBody(tc.requestJSON)
			resp, err := req.Send()
			assert.NoError(t, err, "error making HTTP request")
			assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
			assert.Equal(t, tc.expectedResponseJSON, string(resp.Body()))

		})
	}
}
