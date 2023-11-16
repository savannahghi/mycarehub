package gorm

import (
	"fmt"
	"os"

	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"

	// responsible for providing methods that are by gorm to connect to cloud SQL
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	_ "github.com/lib/pq"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	// DBHost ..
	DBHost = "POSTGRES_HOST"
	// DBPort ...
	DBPort = "POSTGRES_PORT"
	// DBUser ...
	DBUser = "POSTGRES_USER"
	// DBPASSWORD ...
	DBPASSWORD = "POSTGRES_PASSWORD"
	// DBName ...
	DBName = "POSTGRES_DB"
)

type connectionConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// PGInstance box for postgres client. We use this instead of a global variable
type PGInstance struct {
	DB *gorm.DB
}

// NewPGInstance creates a new instance of postgres client
func NewPGInstance() (*PGInstance, error) {
	db := startDatabase()
	if db == nil {
		return nil, fmt.Errorf("failed to start database: %v", db)
	}
	pg := &PGInstance{DB: db}

	return pg, nil
}

// startDatabase ...
func startDatabase() *gorm.DB {
	config := connectionConfig{
		host:     serverutils.MustGetEnvVar(DBHost),
		port:     serverutils.MustGetEnvVar(DBPort),
		user:     serverutils.MustGetEnvVar(DBUser),
		password: serverutils.MustGetEnvVar(DBPASSWORD),
		dbname:   serverutils.MustGetEnvVar(DBName),
	}
	return boot(config)
}

func boot(cfg connectionConfig) *gorm.DB {
	var err error
	var db *gorm.DB

	connString := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v", cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbname)
	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})

	if db == nil || err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// Check that we can connect to the DB. This allows us to know early on if
	// we have the correct credentials and connection settings. As a bonus, we
	// also get a more descriptive error message and exit code.
	if connection, err := db.DB(); err != nil {
		log.Errorf("unable to get a connection to the database: %s", err)
		os.Exit(1)
	} else {
		if err = connection.Ping(); err != nil {
			log.Errorf("unable to ping the database: %s", err)
			os.Exit(1)
		}
	}

	return db
}
