package metrics_test

import (
	"context"
	"testing"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/metrics"
)

func TestUsecaseMetricsImpl_CollectMetric(t *testing.T) {
	type args struct {
		ctx   context.Context
		input domain.Metric
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: record a metric",
			args: args{
				ctx: context.Background(),
				input: domain.Metric{
					Type: enums.MetricTypeContent,
					Event: map[string]interface{}{
						"contentID": 10,
						"duration":  time.Since(time.Now()),
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: invalid metric",
			args: args{
				ctx: context.Background(),
				input: domain.Metric{
					Type: "INVALID",
					Event: map[string]interface{}{
						"contentID": 10,
						"duration":  time.Since(time.Now()),
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			m := metrics.NewUsecaseMetricsImpl(fakeDB)

			got, err := m.CollectMetric(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseMetricsImpl.CollectMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseMetricsImpl.CollectMetric() = %v, want %v", got, tt.want)
			}
		})
	}
}
