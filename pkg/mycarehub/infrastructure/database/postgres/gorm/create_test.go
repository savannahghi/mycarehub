package gorm_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestPGInstance_SavePin(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		pinData *gorm.PINData
	}

	id := uuid.New().String()
	pinPayload := &gorm.PINData{
		PINDataID: &id,
		UserID:    id,
		HashedPIN: "1234",
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   feedlib.FlavourConsumer,
		Salt:      "1234",
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save pin",
			args: args{
				ctx:     ctx,
				pinData: pinPayload,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save pin",
			args: args{
				ctx:     ctx,
				pinData: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SavePin(tt.args.ctx, tt.args.pinData)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SavePin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SavePin() = %v, want %v", got, tt.want)
			}
		})
	}
}
