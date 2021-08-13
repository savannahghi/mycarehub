package domain

import (
	"time"
)

var (
	// TimeLocation ...
	TimeLocation, _ = time.LoadLocation("Africa/Nairobi")

	// TimeFormatStr date time string format
	TimeFormatStr = "2006-01-02T15:04:05+03:00"

	// Repo the env to identify which repo to use
	Repo = "REPOSITORY"

	//FirebaseRepository is the value of the env when using firebase
	FirebaseRepository = "firebase"

	//PostgresRepository is the value of the env when using postgres
	PostgresRepository = "postgres"
)
