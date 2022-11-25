package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func Test_createMapUser(t *testing.T) {

	username := ksuid.New().String()
	programID := uuid.New().String()

	type args struct {
		userObject *gorm.User
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
	}{
		{
			name: "Happy Case",
			args: args{
				&gorm.User{
					Username:         username,
					CurrentProgramID: programID,
				},
			},
			wantNil: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createMapUser(tt.args.userObject)
			if !tt.wantNil && got == nil {
				t.Errorf("wanted a value but got: %v", got)
			}
		})
	}
}
