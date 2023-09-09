package metrics

import "testing"

func TestCounterMetric_GetName(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}
			if got := counterMetric.GetName(); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCounterMetric_UpdateValue(t *testing.T) {
	type fields struct {
		CMetric Metric
		value   uint64
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
			counterMetric := &CounterMetric{
				CMetric: tt.fields.CMetric,
				value:   tt.fields.value,
			}
			if err := counterMetric.UpdateValue(tt.args.newValue); (err != nil) != tt.wantErr {
				t.Errorf("UpdateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
