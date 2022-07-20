package gorm

import (
	"testing"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/savannahghi/serverutils"
	"github.com/tj/assert"
	"gorm.io/gorm"
)

func Test_boot(t *testing.T) {
	type args struct {
		cfg connectionConfig
	}
	tests := []struct {
		name string
		args args
		want *gorm.DB
	}{
		{
			name: "Happy case",
			args: args{
				cfg: connectionConfig{
					host:            serverutils.MustGetEnvVar(DBHost),
					port:            serverutils.MustGetEnvVar(DBPort),
					user:            serverutils.MustGetEnvVar(DBUser),
					password:        serverutils.MustGetEnvVar(DBPASSWORD),
					dbname:          serverutils.MustGetEnvVar(DBName),
					project:         serverutils.MustGetEnvVar(GoogleProject),
					region:          serverutils.MustGetEnvVar(DatabaseRegion),
					instance:        serverutils.MustGetEnvVar(DatabasesInstance),
					asCloudInstance: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boot(tt.args.cfg)
			assert.NotNil(t, got)
		})
	}
}
