package metrics

import "testing"

func TestGaugeMetric_GetName(t *testing.T) {
	type fields struct {
		GMetric Metric
		Value   float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "Simple Test",
			fields: struct {
				GMetric Metric
				Value   float64
			}{GMetric: struct {
				Name  string
				Mtype MetricType
			}{Name: "Test", Mtype: Gauge}, Value: 0},
			want: "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				GMetric: tt.fields.GMetric,
				Value:   tt.fields.Value,
			}
			if got := gaugeMetric.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGaugeMetric_UpdateValue(t *testing.T) {
	type fields struct {
		GMetric Metric
		Value   float64
	}
	type args struct {
		newValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gaugeMetric := &GaugeMetric{
				GMetric: tt.fields.GMetric,
				Value:   tt.fields.Value,
			}
			if err := gaugeMetric.UpdateValue(tt.args.newValue); (err != nil) != tt.wantErr {
				t.Errorf("UpdateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
