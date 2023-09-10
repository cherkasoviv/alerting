package handlers

import (
	"alerting/internal/mstorage"
	"net/http"
	"testing"
)

func TestMetricHandler_UpdateRequest(t *testing.T) {
	type fields struct {
		Storage *mstorage.MetricStorage
	}
	type args struct {
		res http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MetricHandler{
				Storage: tt.fields.Storage,
			}
			m.UpdateRequest(tt.args.res, tt.args.req)
		})
	}
}
