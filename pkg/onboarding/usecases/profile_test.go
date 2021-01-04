package usecases

import (
	"context"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

func TestMaskPhoneNumbers(t *testing.T) {

	ctx := context.Background()
	repository, err := database.NewFirebaseRepository(ctx)
	if err != nil {
		t.Errorf("failed to create a repository")
		return
	}

	profileUseCase := NewProfileUseCase(repository)

	type args struct {
		phones []string
	}

	tests := []struct {
		name string
		arg  args
		want []string
	}{
		{
			name: "valid case",
			arg: args{
				phones: []string{"+254789874267"},
			},
			want: []string{"+254789***267"},
		},
		{
			name: "valid case < 10 digits",
			arg: args{
				phones: []string{"+2547898742"},
			},
			want: []string{"+2547***742"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maskedPhone := profileUseCase.MaskPhoneNumbers(tt.arg.phones)
			if len(maskedPhone) != len(tt.want) {
				t.Errorf("returned masked phone number not the expected one, wanted: %v got: %v", tt.want, maskedPhone)
				return
			}

			for i, number := range maskedPhone {
				if tt.want[i] != number {
					t.Errorf("wanted: %v, got: %v", tt.want[i], number)
					return
				}
			}
		})
	}
}
