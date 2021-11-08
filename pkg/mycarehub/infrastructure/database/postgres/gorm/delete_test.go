package gorm_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_DeleteFacility(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()
	name := ksuid.New().String()
	code := uuid.New().String()
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &gorm.Facility{
		FacilityID:  &ID,
		Name:        name,
		Code:        code,
		Active:      strconv.FormatBool(true),
		County:      county,
		Description: description,
	}

	facility, err := testingDB.GetOrCreateFacility(ctx, facility)
	if err != nil {
		t.Errorf("failed to create test facility")
		return
	}

	type args struct {
		ctx     context.Context
		mflcode string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully delete facility",
			args: args{
				ctx:     ctx,
				mflcode: facility.Code,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to delete facility",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.DeleteFacility(tt.args.ctx, tt.args.mflcode)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.DeleteFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}
