package admin_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/usecases/admin"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("ROOT_COLLECTION_SUFFIX", "testing_ci")

	log.Printf("Running tests ...")

	code := m.Run()

	log.Printf("Tearing tests down ...")

	os.Exit(code)
}

func TestService_CheckPreconditions(t *testing.T) {
	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	goodService := admin.NewService(ext)
	assert.NotNil(t, goodService)
	goodService.CheckPreconditions() // no panic

	badService := admin.Service{}
	assert.Panics(t, func() {
		badService.CheckPreconditions()
	})
}

func cleanup(ctx context.Context, s *admin.Service, t *testing.T) {
	got, _ := s.ListMicroservices(ctx)

	for _, ms := range got {
		_, err := s.DeregisterMicroservice(ctx, ms.ID)
		assert.Nil(t, err)
	}
}

func TestService_RegisterMicroservice(t *testing.T) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	gofakeit.Seed(r1.Int63()) // seed the pseudo-random generator

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	s := admin.NewService(ext)

	// create authenticated context
	ctx := firebasetools.GetAuthenticatedContext(t)
	cleanup(ctx, s, t)

	type args struct {
		ctx   context.Context
		input domain.Microservice
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "good case",
			args: args{
				ctx: ctx,
				input: domain.Microservice{
					Name:        uuid.New().String(),
					Description: uuid.New().String(),
					URL:         "https://profile-staging.healthcloud.co.ke/graphql",
				},
			},
			wantErr: false,
		},

		{
			name: "invalid endpoint",
			args: args{
				ctx: context.Background(),
				input: domain.Microservice{
					Name:        uuid.New().String(),
					Description: uuid.New().String(),
					URL:         "https://profile-staging.healthcloud.co.ke/gra",
				},
			},
			wantErr: true,
		},

		{
			name: "insecure endpoint",
			args: args{
				ctx: context.Background(),
				input: domain.Microservice{
					Name:        uuid.New().String(),
					Description: uuid.New().String(),
					URL:         "http://profile-staging.healthcloud.co.ke/health",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &firebasetools.FirebaseClient{}
			ext := extension.NewBaseExtensionImpl(fc)
			s := admin.NewService(ext)

			got, err := s.RegisterMicroservice(tt.args.ctx, tt.args.input)
			// an error occurred yet it was not expected
			if tt.wantErr == false && err != nil {
				t.Errorf("Service.RegisterMicroservice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// we expect an error
			if tt.wantErr {
				assert.Nil(t, got)

				// re-register...there should be an error
				repeatReg, err := s.RegisterMicroservice(tt.args.ctx, tt.args.input)
				assert.NotNil(t, err)
				assert.Nil(t, repeatReg)
			}
		})
	}

	cleanup(ctx, s, t)
}

func TestService_ListMicroservices(t *testing.T) {
	// make at least one micro-service
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	gofakeit.Seed(r1.Int63()) // seed the pseudo-random

	inp := domain.Microservice{
		Name:        uuid.New().String(),
		Description: uuid.New().String(),
		URL:         "https://profile-staging.healthcloud.co.ke/graphql",
	}

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	srv := admin.NewService(ext)
	ctx := firebasetools.GetAuthenticatedContext(t)

	cleanup(ctx, srv, t)

	_, err := srv.RegisterMicroservice(ctx, inp)
	assert.Nil(t, err)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &firebasetools.FirebaseClient{}
			ext := extension.NewBaseExtensionImpl(fc)
			s := admin.NewService(ext)
			got, err := s.ListMicroservices(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.ListMicroservices() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, len(got), 1)

				for _, ms := range got {
					_, err = s.DeregisterMicroservice(tt.args.ctx, ms.ID)
					assert.Nil(t, err)
				}
			}
		})
	}

	cleanup(ctx, srv, t)
}

func TestService_FindMicroserviceByID(t *testing.T) {
	// make at least one micro-service
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	gofakeit.Seed(r1.Int63()) // seed the pseudo-random

	inp := domain.Microservice{
		Name:        uuid.New().String(),
		Description: uuid.New().String(),
		URL:         "https://profile-staging.healthcloud.co.ke/graphql",
	}

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	srv := admin.NewService(ext)
	ctx := firebasetools.GetAuthenticatedContext(t)

	cleanup(ctx, srv, t)

	microservice, err := srv.RegisterMicroservice(ctx, inp)
	assert.Nil(t, err)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "service that exists",
			args: args{
				ctx: ctx,
				id:  microservice.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &firebasetools.FirebaseClient{}
			ext := extension.NewBaseExtensionImpl(fc)
			s := admin.NewService(ext)
			got, err := s.FindMicroserviceByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.FindMicroserviceByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.NotNil(t, got)
			}
		})
	}

	cleanup(ctx, srv, t)
}

func TestService_CheckHealthEndpoint(t *testing.T) {
	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	s := admin.NewService(ext)
	ctx := context.Background()

	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "valid_case",
			args: "https://profile-staging.healthcloud.co.ke/health",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := s.CheckHealthEndpoint(ctx, tt.args); got != tt.want {
				t.Errorf("CheckHealthEndpoint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeregisterAllServices(t *testing.T) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	gofakeit.Seed(r1.Int63()) // seed the pseudo-random generator

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	s := admin.NewService(ext)
	ctx := firebasetools.GetAuthenticatedContext(t)
	cleanup(ctx, s, t)

	// register a few services
	services := []domain.Microservice{
		{
			Name:        "engagement staging",
			Description: uuid.New().String(),
			URL:         "https://profile-staging.healthcloud.co.ke/graphql",
		},
	}

	for _, srv := range services {
		_, err := s.RegisterMicroservice(ctx, srv)
		if err != nil {
			t.Fatalf("failed to register test services : %v", err)
		}
	}

	srvs1, err := s.ListMicroservices(ctx)
	if err != nil {
		t.Fatalf("failed to list test services (srvs1): %v", err)
	}

	// expect to be equal to 2
	assert.Equal(t, len(srvs1), 1, "Expected 1 service only")

	_, errd := s.DeregisterAllMicroservices(ctx)
	if errd != nil {
		t.Fatalf("failed to deregister all test services : %v", err)
	}

	srvs2, err := s.ListMicroservices(ctx)
	if err != nil {
		t.Fatalf("failed to list test services (srvs2): %v", err)
	}

	assert.Equal(t, len(srvs2), 0, fmt.Sprintf("Expected 0 services, Got %v", len(srvs2)))

	cleanup(ctx, s, t)

}

func TestDeregisterServiceWithID(t *testing.T) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	gofakeit.Seed(r1.Int63()) // seed the pseudo-random generator

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	s := admin.NewService(ext)
	ctx := firebasetools.GetAuthenticatedContext(t)
	cleanup(ctx, s, t)

	// register a few services
	services := []domain.Microservice{
		{
			Name:        "engagement staging",
			Description: uuid.New().String(),
			URL:         "https://engagement-staging.healthcloud.co.ke/graphql",
		},
		{
			Name:        "profile staging",
			Description: uuid.New().String(),
			URL:         "https://profile-staging.healthcloud.co.ke/graphql",
		},
	}

	for _, srv := range services {
		_, err := s.RegisterMicroservice(ctx, srv)
		if err != nil {
			t.Fatalf("failed to register test services : %v", err)
		}
	}

	srvs, err := s.ListMicroservices(ctx)
	if err != nil {
		t.Fatalf("failed to list test services (srvs1): %v", err)
	}

	// expect to be equal to 2
	assert.Equal(t, len(srvs), 2, "Expected 2 services only")

	_, errd := s.DeregisterMicroservice(ctx, srvs[0].ID)
	if errd != nil {
		t.Fatalf("failed to deregister test service %v : %v", srvs[0].Name, err)
	}

	srvs1, err := s.ListMicroservices(ctx)
	if err != nil {
		t.Fatalf("failed to list test services (srvs1): %v", err)
	}

	// expect to be equal to 1
	assert.Equal(t, len(srvs1), 1, "Expected 1 service only")

	_, errd2 := s.DeregisterMicroservice(ctx, srvs[1].ID)
	if errd2 != nil {
		t.Fatalf("failed to deregister service %v : %v", srvs[1].Name, err)
	}

	srvs2, err := s.ListMicroservices(ctx)
	if err != nil {
		t.Fatalf("failed to list test services (srvs2): %v", err)
	}

	// expect to be equal to 0
	assert.Equal(t, len(srvs2), 0, "Expected 0 services")

	cleanup(ctx, s, t)

}

func TestService_PollMicroservicesStatus(t *testing.T) {

	fc := &firebasetools.FirebaseClient{}
	ext := extension.NewBaseExtensionImpl(fc)
	s := admin.NewService(ext)

	input := domain.Microservice{

		Name:        uuid.New().String(),
		Description: uuid.New().String(),
		URL:         "https://profile-staging.healthcloud.co.ke/graphql",
	}
	ctx := firebasetools.GetAuthenticatedContext(t)
	cleanup(ctx, s, t)

	_, err := s.RegisterMicroservice(ctx, input)
	assert.Nil(t, err)

	services, _ := s.ListMicroservices(ctx)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.MicroserviceStatus
		wantErr bool
	}{
		{
			name: "poll services",
			args: args{
				ctx: ctx,
			},
			want: []*domain.MicroserviceStatus{
				{
					Service: services[0],
					Active:  true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.PollMicroservicesStatus(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.PollMicroservicesStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.PollMicroservicesStatus() = %v, want %v", got, tt.want)
			}

			if !tt.wantErr {
				assert.Equal(t, len(got), 1)

				for _, ms := range got {
					_, err = s.DeregisterMicroservice(tt.args.ctx, ms.Service.ID)
					assert.Nil(t, err)
				}
			}
		})
	}

	cleanup(ctx, s, t)
}
