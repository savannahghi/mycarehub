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

//WelcomeMessage is the default message formart for sending temporary PIN to users
var WelcomeMessage = "Hi %s, welcome to Be.Well. Please use this One Time PIN: %s to log in using your phone number. You will be prompted to set a new PIN on login."
