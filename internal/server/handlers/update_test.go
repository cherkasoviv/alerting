package handlers

import (
	"alerting/internal/config"
	mstorage "alerting/internal/mstorage"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func TestMetricHandlerUpdateRequest(t *testing.T) {
	cfg := config.ServerConfig{
		Host:            "",
		StoreInterval:   0,
		FileStoragePath: "",
		NeedToRestore:   false,
	}
	storage := mstorage.Initialize(&cfg)
	updateHandler := NewUpdateHandler(storage)
	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Route("/{metricType}/{metricName}/{metricValue}", func(r chi.Router) {
			r.Post("/", updateHandler.CreateOrUpdateFromURLPath())
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

func TestUpdateHandlerCreateOrUpdateFromJSON(t *testing.T) {
	cfg := config.ServerConfig{
		Host:            "",
		StoreInterval:   0,
		FileStoragePath: "",
		NeedToRestore:   false,
	}
	storage := mstorage.Initialize(&cfg)
	updateHandler := NewUpdateHandler(storage)

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/", updateHandler.CreateOrUpdateFromJSON())
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
		{name: "success gauge", method: http.MethodPost, requestURL: "/update", expectedCode: http.StatusOK,
			requestJSON:          "{\"id\":\"testGauge\" , \"type\": \"gauge\",\"value\":1}\n}",
			expectedResponseJSON: "{\"id\":\"testGauge\",\"type\":\"gauge\",\"value\":1}\n"},
		{name: "success counter", method: http.MethodPost, requestURL: "/update", expectedCode: http.StatusOK,
			requestJSON:          "{\"id\":\"testCounter\" , \"type\": \"counter\",\"delta\":1}\n}",
			expectedResponseJSON: "{\"id\":\"testCounter\",\"type\":\"counter\",\"delta\":1}\n"},
		{name: "error json", method: http.MethodPost, requestURL: "/update", expectedCode: http.StatusBadRequest,
			requestJSON:          "{\"name\":\"test\" ,\"type\":\"gauge\"}",
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

func BenchmarkUpdateHandlerCreateOrUpdateFromJSON(b *testing.B) {
	cfg := config.ServerConfig{
		Host:            "",
		StoreInterval:   0,
		FileStoragePath: "",
		NeedToRestore:   false,
	}
	storage := mstorage.Initialize(&cfg)
	updateHandler := NewUpdateHandler(storage)

	r := chi.NewRouter()
	r.Route("/update", func(r chi.Router) {
		r.Post("/", updateHandler.CreateOrUpdateFromJSON())
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	testCase := struct {
		method     string
		requestURL string

		requestJSON string
	}{method: http.MethodPost, requestURL: "/update",
		requestJSON: "{\"id\":\"testGauge\" , \"type\": \"gauge\",\"value\":1}\n}",
	}

	for i := 0; i < b.N; i++ {

		req := resty.New().R()
		req.Method = testCase.method
		req.URL = ts.URL + testCase.requestURL
		req.SetBody(testCase.requestJSON)
		req.Send()

	}
}
