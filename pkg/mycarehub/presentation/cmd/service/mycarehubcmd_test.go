package service_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/cmd/service"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	appointmentMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/appointments/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	communitiesMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities/mock"
	contentMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content/mock"
	facilityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
	feedbackMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback/mock"
	healthdiaryMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary/mock"
	metricsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/metrics/mock"
	notificationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
	oauthMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth/mock"
	organisationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	programsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/pubsub/mock"
	questionnairesMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires/mock"
	securityquestionsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions/mock"
	servicerequestMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	surveysMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys/mock"
	termsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
	userMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func TestMyCareHubCmdInterfacesImpl_CreateSuperUser(t *testing.T) {
	type createsuperuserInput struct {
		organisationIndex string
		programIndex      string
		username          string
		firstName         string
		lastName          string
		birthYear         string
		birthMonth        string
		birthDay          string
		gender            string
		phone             string
		sendInvite        string
		idNumber          string
		staffNumber       string
	}
	type args struct {
		ctx   context.Context
		input createsuperuserInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: create superuser",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid organisation selection",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "40",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty organisation list",
			args: args{
				ctx: context.Background(),
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid program selection",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "40",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty program list",
			args: args{
				ctx: context.Background(),
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: username not alphanumeric",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "inv@lid",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing first name",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing last name",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing year",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid year",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "1000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing month",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid month",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "20",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing day",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid day",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "50",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing gender",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid gender",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.Gender("invalid").String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing phone",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid phone",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "invalid",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing sendInvite",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid sendInvite",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "invalid",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing id number",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing staff number",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to check if superuser exists",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: superuser exists",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get program facilities",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: program facilities not found",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create superuser",
			args: args{
				ctx: nil,
				input: createsuperuserInput{
					organisationIndex: "0",
					programIndex:      "0",
					username:          "username",
					firstName:         gofakeit.Name(),
					lastName:          gofakeit.Name(),
					birthYear:         "2000",
					birthMonth:        "1",
					birthDay:          "1",
					gender:            enumutils.GenderMale.String(),
					phone:             "0999999999",
					sendInvite:        "yes",
					idNumber:          "328392893082903",
					staffNumber:       "st323232",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase, organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)

			if tt.name == "Sad Case: failed to check if superuser exists" {
				userUsecase.MockCheckSuperUserExistsFn = func(ctx context.Context) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}

			}
			if tt.name == "Sad Case: superuser exists" {
				userUsecase.MockCheckSuperUserExistsFn = func(ctx context.Context) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad Case: failed to get programs" {
				programsUsecase.MockListProgramsFn = func(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: programs not found" {
				programsUsecase.MockListProgramsFn = func(ctx context.Context, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
					return &domain.ProgramPage{
						Pagination: domain.Pagination{},
						Programs:   []*domain.Program{},
					}, nil
				}
			}

			if tt.name == "Sad Case: failed to get program facilities" {
				programsUsecase.MockGetProgramFacilitiesFn = func(ctx context.Context, programID string) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: program facilities not found" {
				programsUsecase.MockGetProgramFacilitiesFn = func(ctx context.Context, programID string) ([]*domain.Facility, error) {
					return []*domain.Facility{}, nil
				}
			}

			if tt.name == "Sad Case: failed to create superuser" {
				userUsecase.MockCreateSuperUserFn = func(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: empty organisation list" {
				organisationUsecase.MockListOrganisationsFn = func(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
					return nil, nil
				}
			}
			if tt.name == "Sad Case: empty program list" {
				programsUsecase.MockListOrganisationProgramsFn = func(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
					return nil, nil
				}
			}

			stdoutString := fmt.Sprintf("%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n",
				tt.args.input.organisationIndex,
				tt.args.input.programIndex,
				tt.args.input.username,
				tt.args.input.firstName,
				tt.args.input.lastName,
				tt.args.input.birthYear,
				tt.args.input.birthMonth,
				tt.args.input.birthDay,
				tt.args.input.gender,
				tt.args.input.phone,
				tt.args.input.sendInvite,
				tt.args.input.idNumber,
				tt.args.input.staffNumber,
			)
			input := bytes.NewBufferString(stdoutString)
			m := service.NewMyCareHubCmdInterfaces(*usecases)
			if err := m.CreateSuperUser(tt.args.ctx, input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.CreateSuperUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LoadFacilities(t *testing.T) {
	type args struct {
		ctx              context.Context
		absoluteFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Load Facilities",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/valid.csv",
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid path",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "invalid/test.csv",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/invalidPhone.csv",
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing field value",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/missingFieldValue.csv",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid phone",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/invalidPhone.csv",
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to create facility",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/valid.csv",
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to publish facilities to CMS",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/facility/valid.csv",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad case: failed to create facility" {
				facilityUseCase.MockCreateFacilitiesFn = func(ctx context.Context, facilities []*domain.Facility) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to publish facilities to CMS" {
				facilityUseCase.MockPublishFacilitiesToCMSFn = func(ctx context.Context, facilities []*domain.Facility) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if err := m.LoadFacilities(tt.args.ctx, tt.args.absoluteFilePath); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LoadFacilities() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LoadOrganisation(t *testing.T) {
	type args struct {
		ctx              context.Context
		organisationPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: load organisation",
			args: args{
				ctx:              context.Background(),
				organisationPath: "testData/organisation/valid.json",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid json field to map to organisation",
			args: args{
				ctx:              context.Background(),
				organisationPath: "testData/organisation/invalidField.json",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid json file to map to organisation",
			args: args{
				ctx:              context.Background(),
				organisationPath: "testData/organisation/invalidJson",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid json file path for organisation",
			args: args{
				ctx:              context.Background(),
				organisationPath: "invalidPath",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create organisation",
			args: args{
				ctx:              context.Background(),
				organisationPath: "testData/organisation/valid.json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communitiesUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecase := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase, organisationUsecase, pubSubUseCase, communitiesUsecase, oauthUsecase,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad Case: failed to create organisation" {
				organisationUsecase.MockCreateOrganisationFn = func(ctx context.Context, input dto.OrganisationInput, programInput []*dto.ProgramInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if err := m.LoadOrganisation(tt.args.ctx, tt.args.organisationPath); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LoadOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LoadProgram(t *testing.T) {
	type loadProgramInput struct {
		organisationIndex string
	}
	type args struct {
		ctx         context.Context
		programPath string
		input       loadProgramInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: load program",
			args: args{
				ctx:         context.Background(),
				programPath: "testData/program/valid.json",
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid json field to map to program",
			args: args{
				ctx:         context.Background(),
				programPath: "testData/program/invalidField.json",
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid organisation index",
			args: args{
				ctx:         context.Background(),
				programPath: "testData/program/valid.json",
				input: loadProgramInput{
					organisationIndex: "40",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid json file to map to program",
			args: args{
				ctx:         context.Background(),
				programPath: "testData/program/invalidJson",
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid json file path for program",
			args: args{
				ctx:         context.Background(),
				programPath: "invalidPath",
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create Program",
			args: args{
				ctx:         context.Background(),
				programPath: "testData/program/valid.json",
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty organisation list",
			args: args{
				ctx: context.Background(),
				input: loadProgramInput{
					organisationIndex: "0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad Case: failed to create Program" {
				programsUsecase.MockCreateProgramFn = func(ctx context.Context, input *dto.ProgramInput) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: empty organisation list" {
				organisationUsecase.MockListOrganisationsFn = func(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
					return nil, nil
				}
			}

			stdoutString := fmt.Sprintf("%v\n",
				tt.args.input.organisationIndex,
			)
			input := bytes.NewBufferString(stdoutString)

			if err := m.LoadProgram(tt.args.ctx, tt.args.programPath, input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LoadOrganisatioAndProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LinkFacilityToProgram(t *testing.T) {
	type linkFacilityToProgramInput struct {
		organisationIndex string
		programIndex      string
		facilityIndex     string
	}
	type args struct {
		ctx   context.Context
		input linkFacilityToProgramInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: link facility to program",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid organisation selection",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "40",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid program selection",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "40",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid facility selection",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "40",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to link facility to program",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty organisation list",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty program list",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: empty facility list",
			args: args{
				ctx: context.Background(),
				input: linkFacilityToProgramInput{
					organisationIndex: "0",
					programIndex:      "0",
					facilityIndex:     "0",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad Case: failed to link facility to program" {
				facilityUseCase.MockAddFacilityToProgramFn = func(ctx context.Context, facilityIDs []string, programID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad Case: empty organisation list" {
				organisationUsecase.MockListOrganisationsFn = func(ctx context.Context, paginationInput *dto.PaginationsInput) (*dto.OrganisationOutputPage, error) {
					return nil, nil
				}
			}
			if tt.name == "Sad Case: empty program list" {
				programsUsecase.MockListOrganisationProgramsFn = func(ctx context.Context, organisationID string, paginationsInput *dto.PaginationsInput) (*domain.ProgramPage, error) {
					return nil, nil
				}
			}
			if tt.name == "Sad Case: empty facility list" {
				facilityUseCase.MockListFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, nil
				}
			}

			stdoutString := fmt.Sprintf("%v\n%v\n%v\n",
				tt.args.input.organisationIndex,
				tt.args.input.programIndex,
				tt.args.input.facilityIndex,
			)
			input := bytes.NewBufferString(stdoutString)
			if err := m.LinkFacilityToProgram(tt.args.ctx, input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LinkFacilityToProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LoadSecurityQuestions(t *testing.T) {
	type args struct {
		ctx              context.Context
		absoluteFilePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: load security questions",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/securityquestions/valid.json",
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid flavour",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/securityquestions/invalidFlavour.json",
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid flavour",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/securityquestions/invalidResponseType.json",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create security questions",
			args: args{
				ctx:              context.Background(),
				absoluteFilePath: "testData/securityquestions/valid.json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad Case: failed to create security questions" {
				securityQuestionsUsecase.MockCreateSecurityQuestionsFn = func(ctx context.Context, securityQuestions []*domain.SecurityQuestion) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if err := m.LoadSecurityQuestions(tt.args.ctx, tt.args.absoluteFilePath); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LoadSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubCmdInterfacesImpl_LoadTermsOfService(t *testing.T) {
	type termsInput struct {
		path  string
		years string
	}
	type args struct {
		ctx   context.Context
		input termsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: load terms of service",
			args: args{
				ctx: nil,
				input: termsInput{
					path:  "testData/terms/terms.txt",
					years: "5",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid path",
			args: args{
				ctx: nil,
				input: termsInput{
					path:  "invalid",
					years: "5",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid year input",
			args: args{
				ctx: nil,
				input: termsInput{
					path:  "testData/terms/terms.txt",
					years: "five",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create terms of service",
			args: args{
				ctx: nil,
				input: termsInput{
					path:  "testData/terms/terms.txt",
					years: "5",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			m := service.NewMyCareHubCmdInterfaces(*usecases)

			if tt.name == "Sad Case: failed to create terms of service" {
				termsUsecase.MockCreateTermsOfServiceFn = func(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			stdoutString := fmt.Sprintf("%v\n%v\n",
				tt.args.input.path,
				tt.args.input.years,
			)
			input := bytes.NewBufferString(stdoutString)
			if err := m.LoadTermsOfService(tt.args.ctx, input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.LoadTermsOfService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
