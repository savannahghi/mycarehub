package gorm_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
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

func TestPGInstance_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	ID := uuid.New().String()
	name := gofakeit.Name()
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

	type args struct {
		ctx      context.Context
		facility *gorm.Facility
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get or create facility",
			args: args{
				ctx:      ctx,
				facility: facility,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail tp get or create facility",
			args: args{
				ctx:      ctx,
				facility: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestPGInstance_RegisterClient(t *testing.T) {
	ctx := context.Background()
	User := &gorm.User{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Username:    gofakeit.Username(),
		MiddleName:  gofakeit.Name(),
		DisplayName: gofakeit.BeerAlcohol(),
		Gender:      enumutils.GenderMale,
	}
	Client := &gorm.ClientProfile{
		ClientType: enums.ClientTypeOvc,
	}
	type args struct {
		ctx           context.Context
		user          *gorm.User
		clientProfile *gorm.ClientProfile
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register client",
			args: args{
				ctx:           ctx,
				user:          User,
				clientProfile: Client,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RegisterClient(tt.args.ctx, tt.args.user, tt.args.clientProfile)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
