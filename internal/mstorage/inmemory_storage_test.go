package mstorage

import (
	metric "alerting/internal/metrics"
	"reflect"
	"testing"
)

func TestInMemorystorage_FindAllMetrics(t *testing.T) {
	type fields struct {
		storage *map[string]metric.AbstractMetric
	}
	tests := []struct {
		name    string
		fields  fields
		want    *map[string]metric.AbstractMetric
		wantErr bool
	}{
		{
			name: "Empty storage",
			fields: struct {
				storage *map[string]metric.AbstractMetric
			}{storage: &map[string]metric.AbstractMetric{}},
			want:    &map[string]metric.AbstractMetric{},
			wantErr: false,
		},
		{
			name: "Nonempty storage",
			fields: struct {
				storage *map[string]metric.AbstractMetric
			}{storage: &map[string]metric.AbstractMetric{
				"testMetric": &metric.GaugeMetric{
					GMetric: struct {
						Name  string
						Mtype metric.MetricType
					}{Name: "testMetric", Mtype: metric.Gauge},
					Value: 0,
				},
			}},
			want: &map[string]metric.AbstractMetric{
				"testMetric": &metric.GaugeMetric{
					GMetric: struct {
						Name  string
						Mtype metric.MetricType
					}{Name: "testMetric", Mtype: metric.Gauge},
					Value: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := InMemoryStorage{
				storage: tt.fields.storage,
			}
			got, err := st.FindAllMetrics()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, *tt.want) {
				t.Errorf("FindAllMetrics() got = %v, want %v", got, tt.want)
			}
		})
	}
}
