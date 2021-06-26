package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestUSSDUseCaseImpl_GenerateUSSD(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	phone := "+254721026491"
	invalidPhone := ""
	validInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
	}
	invalidInput := dto.SessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &invalidPhone,
	}
	aplhaNumbericphone := "+254-not-valid-123"
	alphanumericPhoneInput := dto.SessionDetails{
		SessionID:   "123455678",
		PhoneNumber: &aplhaNumbericphone,
	}
	level1ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname",
	}
	level1InvalidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*",
	}
	levelInvalidtextInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*1",
	}

	level2ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname",
	}
	level2InvalidtextInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*1",
	}
	level3ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021",
	}
	level3InvalidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25234",
	}
	level4ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*1234",
	}
	level4InvalidPINInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*12",
	}
	// level5ValidInput := dto.SessionDetails{
	// 	SessionID:   "12345678",
	// 	PhoneNumber: &phone,
	// 	Text:        "1*firstname*lastname*25062021*1234*1234",
	// }
	// level5InvalidInput := dto.SessionDetails{
	// 	SessionID:   "12345678",
	// 	PhoneNumber: &phone,
	// 	Text:        "1*firstname*lastname*25062021*1234*4321",
	// }
	level6ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*1234*1234*1",
	}
	level6ValidInputCase2 := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*1234*1234*2",
	}
	// level6ValidInputDefaultCase := dto.SessionDetails{
	// 	SessionID:   "12345678",
	// 	PhoneNumber: &phone,
	// 	Text:        "1*firstname*lastname*25062021*1234*1234*3",
	// }
	level7ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*1234*1234*1*1234",
	}
	level8ValidInput := dto.SessionDetails{
		SessionID:   "12345678",
		PhoneNumber: &phone,
		Text:        "1*firstname*lastname*25062021*1234*1234*1*1234*4321",
	}

	// level9ValidInput := dto.SessionDetails{
	// 	SessionID:   "12345678",
	// 	PhoneNumber: &phone,
	// 	Text:        "1*firstname*lastname*25062021*1234*1234*1*1234*4321*4321",
	// }

	type args struct {
		ctx   context.Context
		input dto.SessionDetails
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "happy:) level 0 success",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want: "CON Welcome to Be.Well \n1. Register",
		},
		{
			name: "sad:( level 0 fail",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want: "failed to add USSD details",
		},
		{
			name: "happy:) valid phone number success",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			want: "CON Welcome to Be.Well \n1. Register",
		},
		{
			name: "sad:( invalid phone number failure",
			args: args{
				ctx:   ctx,
				input: invalidInput,
			},
			want: "2: unable to normalize the msisdn",
		},
		{
			name: "sad:( aplha numeric phone number failure",
			args: args{
				ctx:   ctx,
				input: alphanumericPhoneInput,
			},
			want: "2: unable to normalize the msisdn",
		},
		{
			name: "happy:) level 1 success",
			args: args{
				ctx:   ctx,
				input: level1ValidInput,
			},
			want: "CON Please enter your last name (eg Doe)",
		},
		{
			name: "sad:level 1 empty USSD text input failure",
			args: args{
				ctx:   ctx,
				input: level1InvalidInput,
			},
			want: "CON Invalid name. Please enter a valid name (e.g John)",
		},
		{
			name: "sad:level 1 digit text input failure",
			args: args{
				ctx:   ctx,
				input: levelInvalidtextInput,
			},
			want: "CON Invalid name. Please enter a valid name (e.g John)",
		},
		{
			name: "happy:) level 2 success",
			args: args{
				ctx:   ctx,
				input: level2ValidInput,
			},
			want: "CON Please enter your date of birth in DDMMYYYY format e.g 14031996 for 14th March 1992",
		},
		{
			name: "sad:level 2 digit text input failure",
			args: args{
				ctx:   ctx,
				input: level2InvalidtextInput,
			},
			want: "CON Invalid name. Please enter a valid name (e.g John)",
		},
		{
			name: "happy:) level 3 success",
			args: args{
				ctx:   ctx,
				input: level3ValidInput,
			},
			want: "CON Please enter a 4 digit PIN to secure your account",
		},
		{
			name: "sad:level 3 invalid date format failure",
			args: args{
				ctx:   ctx,
				input: level3InvalidInput,
			},
			want: "CON The date of birth you entered is not valid, please try again in DDMMYYYY format e.g 14031996",
		},
		{
			name: "happy:) level 4 success",
			args: args{
				ctx:   ctx,
				input: level4ValidInput,
			},
			want: "CON Please enter a 4 digit PIN again to confirm",
		},
		{
			name: "sad:) level 4 invalid pin failure",
			args: args{
				ctx:   ctx,
				input: level4InvalidPINInput,
			},
			want: "CON Invalid PIN. Please enter a 4 digit PIN to secure your account",
		},
		// {
		// 	name: "happy:) level 5 success",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: level5ValidInput,
		// 	},
		// 	want: "CON Thanks for signing up for Be.Well \n1. Opt out from marketing messages \n2. Change PIN",
		// },
		// {
		// 	name: "sad:) level 5 pin mismatch failure",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: level5InvalidInput,
		// 	},
		// 	want: "CON PIN mismatch. Please enter a PIN that matches the first PIN",
		// },
		{
			name: "happy:) level 6 case 1 success",
			args: args{
				ctx:   ctx,
				input: level6ValidInput,
			},
			want: "END We have successfully opted you out of marketing messages",
		},
		{
			name: "happy:( level 6 case 2 success",
			args: args{
				ctx:   ctx,
				input: level6ValidInputCase2,
			},
			want: "CON Enter your old PIN to continue",
		},
		// {
		// 	name: "happy:( level 6 default case success",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: level6ValidInputDefaultCase,
		// 	},
		// 	want: "CON Invalid choice. Please try again.\n1.Opt out from marketing messages \n2. Change PIN",
		// },
		{
			name: "happy:) level 7 success",
			args: args{
				ctx:   ctx,
				input: level7ValidInput,
			},
			want: "CON Enter a new four digit PIN",
		},
		{
			name: "sad:( level 7 wrong old pin failure",
			args: args{
				ctx:   ctx,
				input: level7ValidInput,
			},
			want: "CON PIN mismatch. Please enter a PIN that matches your current PIN",
		},
		{
			name: "happy:) level 8 success",
			args: args{
				ctx:   ctx,
				input: level8ValidInput,
			},
			want: "CON Please enter the 4 digit PIN again to confirm",
		},
		// {
		// 	name: "happy:) level 9 success",
		// 	args: args{
		// 		ctx:   ctx,
		// 		input: level9ValidInput,
		// 	},
		// 	want: "END Your PIN was changed successfully",
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) level 0 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 1,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}

			}
			if tt.name == "sad:( level 0 fail" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return fmt.Errorf("failed to add USSD details")
				}
			}
			if tt.name == "happy:) valid phone number success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 1,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721026491"
					return &phone, nil
				}

			}
			if tt.name == "sad:( invalid phone number failure" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "sad:( aplha numeric phone number failure" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("an error occurred")
				}
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "happy:) level 1 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 1,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
			}
			if tt.name == "sad:level 1 empty USSD text input failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "sad:level 1 digit text input failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "happy:) level 2 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 2,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
			}
			if tt.name == "sad:level 2 digit text input failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "sad:) level 2 update session failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 2,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, fmt.Errorf("failed to update session details")
				}
			}
			if tt.name == "happy:) level 3 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 3,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}

			}
			if tt.name == "sad:level 3 invalid date format failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			if tt.name == "sad:) level 3 update session failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 3,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, fmt.Errorf("failed to update session details")
				}
			}
			if tt.name == "happy:) level 4 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 4,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}

			}
			if tt.name == "sad:) level 4 invalid pin failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}

			}
			// if tt.name == "happy:) level 5 success" {
			// 	fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
			// 		return nil
			// 	}
			// 	fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return &domain.USSDLeadDetails{
			// 			Level: 5,
			// 			PIN:   "1234",
			// 		}, nil
			// 	}
			// 	fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return nil, nil
			// 	}
			// 	fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
			// 		return nil, nil
			// 	}
			// fakeRepo.GetOrCreatePhoneNumberUserFn = func(ctx context.Context, phone string) (*dto.CreatedUserResponse, error) {
			// 	phoneNumber := "0715893271"
			// 	return &dto.CreatedUserResponse{
			// 		UID:         "5cf354a2-1d3e-400d-8716-7e2aead29f2c",
			// 		DisplayName: "John Doe",
			// 		Email:       "johndoe@gmail.com",
			// 		PhoneNumber: phoneNumber,
			// 	}, nil
			// }
			// fakeRepo.CreateUserProfileFn = func(ctx context.Context, phoneNumber, uid string) (*base.UserProfile, error) {
			// 	return &base.UserProfile{
			// 			ID: "123",
			// 			VerifiedIdentifiers: []base.VerifiedIdentifier{
			// 				{
			// 					UID:           "125",
			// 					LoginProvider: "Phone",
			// 				},
			// 			},
			// 			PrimaryPhone: &phoneNumber,
			// 		}, nil
			// 	}
			// }

			// if tt.name == "sad:) level 5 pin mismatch failure" {
			// 	fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
			// 		return nil
			// 	}

			// }
			if tt.name == "happy:) level 6 case 1 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 6,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
			}
			if tt.name == "happy:( level 6 case 2 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 6,
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
			}
			// if tt.name == "happy:( level 6 default case success" {
			// 	fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
			// 		return nil
			// 	}
			// 	fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return &domain.USSDLeadDetails{
			// 			Level: 6,
			// 		}, nil
			// 	}
			// 	fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return nil, nil
			// 	}
			// 	fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
			// 		return nil, nil
			// 	}
			// }
			if tt.name == "happy:) level 7 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 7,
						PIN:   "1234",
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
			}

			if tt.name == "sad:( level 7 wrong old pin failure" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 7,
						PIN:   "4321",
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
			}
			if tt.name == "happy:) level 8 success" {
				fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
					return nil
				}
				fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
					return &domain.USSDLeadDetails{
						Level: 8,
						PIN:   "4321",
					}, nil
				}
				fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
					return nil, nil
				}
				fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
					return nil, nil
				}
			}

			// if tt.name == "happy:) level 9 success" {
			// 	fakeRepo.AddAITSessionDetailsFn = func(ctx context.Context, input *dto.SessionDetails) error {
			// 		return nil
			// 	}
			// 	fakeRepo.GetAITSessionDetailsFn = func(ctx context.Context, sessionID string) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return &domain.USSDLeadDetails{
			// 			Level: 9,
			// 			PIN:   "4321",
			// 		}, nil
			// 	}
			// 	fakeRepo.UpdateSessionLevelFn = func(ctx context.Context, sessionID string, level int) (sessionDetails *domain.USSDLeadDetails, err error) {
			// 		return nil, nil
			// 	}
			// 	fakeRepo.UpdateSessionPINFn = func(ctx context.Context, sessionID, pin string) (*domain.USSDLeadDetails, error) {
			// 		return nil, nil
			// 	}
			// }

			resp := i.AITUSSD.GenerateUSSD(tt.args.ctx, &tt.args.input)
			if resp != tt.want {
				t.Errorf("USSDUseCaseImpl.AddUSSDDetails() resp = %v, want %v", resp, tt.want)
				return
			}
		})
	}
}
