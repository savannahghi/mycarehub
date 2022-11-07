package authorization

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/casbin/casbin/v2"
	casbinpgadapter "github.com/cychiuae/casbin-pg-adapter"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/serverutils"
)

func initEnforcer() (*casbin.Enforcer, error) {
	configPath, err := filepath.Abs("./casbin/rbac_model.conf")
	if err != nil {
		return nil, err
	}
	const policyTableName = "casbin_policy"

	db, err := sql.Open("postgres", serverutils.MustGetEnvVar("POSTGRESQL_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	adapter, err := casbinpgadapter.NewAdapter(db, policyTableName)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize casbin pgadapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(configPath, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %v", err)
	}

	return enforcer, nil
}

type Authorization interface {
	LoadPolicy(ctx context.Context)
	Enforce(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
	AddPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
	RemovePolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
	SavePolicy(ctx context.Context) error
	AddGroupingPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
	RemoveGroupingPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
}

type AuthorizationImpl struct {
	enforcer *casbin.Enforcer
}

func NewAuthorizationImpl() Authorization {
	enforcer, err := initEnforcer()
	if err != nil {
		log.Panicf("unable to initialize casbin Enforcer: %s", err)
	}
	return AuthorizationImpl{
		enforcer: enforcer,
	}
}

func (a AuthorizationImpl) LoadPolicy(ctx context.Context) {
	a.enforcer.LoadPolicy()
}
func (a AuthorizationImpl) Enforce(ctx context.Context, subject string, input dto.PermissionInput) (bool, error) {
	return a.enforcer.Enforce(input.OrganizationID, input.ProgramID, subject, input.Object, input.Action)
}
func (a AuthorizationImpl) AddPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error) {
	return a.enforcer.AddPolicy(input.OrganizationID, input.ProgramID, subject, input.Object, input.Action)
}
func (a AuthorizationImpl) RemovePolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error) {
	return a.enforcer.RemovePolicy(input.OrganizationID, input.ProgramID, subject, input.Object, input.Action)
}
func (a AuthorizationImpl) SavePolicy(ctx context.Context) error {
	return a.enforcer.SavePolicy()
}

func (a AuthorizationImpl) AddGroupingPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error) {
	return a.enforcer.AddGroupingPolicy(input.OrganizationID, input.ProgramID, subject, input.Object, input.Action)
}

func (a AuthorizationImpl) RemoveGroupingPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error) {
	return a.enforcer.RemoveGroupingPolicy(input.OrganizationID, input.ProgramID, subject, input.Object, input.Action)
}
