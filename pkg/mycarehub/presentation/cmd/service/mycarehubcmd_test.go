package service_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

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
	organisationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	programsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/pubsub/mock"
	questionnairesMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires/mock"
	screeningtoolsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools/mock"
	securityquestionsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions/mock"
	servicerequestMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	surveysMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys/mock"
	termsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
	userMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func TestMyCareHubCmdInterfacesImpl_CreateSuperUser(t *testing.T) {
	type args struct {
		ctx   context.Context
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: create superuser",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: username not alphanumeric",
			args: args{
				ctx:   nil,
				input: "tes@tuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing first name",
			args: args{
				ctx:   nil,
				input: "testuser\n\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing last name",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\n\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing year",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid year",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n20200\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing month",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid month",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n100\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing day",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid day",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n100\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing gender",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\n\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid gender",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\ninvalid_gender\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing phone",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid phone",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n399939393939393939399393393\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing sendInvite",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\n\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid sendInvite",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\ntrue\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing id number",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing staff number",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: missing facility",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: invalid facility",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n10",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to check if superuser exists",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: superuser exists",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get programs",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: programs not found",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to get program facilities",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: program facilities not found",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
			},
			wantErr: true,
		},
		{
			name: "Sad Case: failed to create superuser",
			args: args{
				ctx:   nil,
				input: "testuser\nfname\nlname\n2020\n01\n01\nmale\n+254999999999\nno\n12121212\n4493943994\n0\n0\n",
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
			communitiesUseCase := communitiesMock.NewCommunityUsecaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			screeningToolsUsecases := screeningtoolsMock.NewScreeningToolsUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			usecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase, communitiesUseCase, screeningToolsUsecases,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase,
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

			input := bytes.NewBufferString(tt.args.input)
			m := service.NewMyCareHubCmdInterfaces(*usecases)
			if err := m.CreateSuperUser(tt.args.ctx, input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubCmdInterfacesImpl.CreateSuperUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
