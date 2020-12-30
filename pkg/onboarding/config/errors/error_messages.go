package errors

const (
	// UserNotFoundErrMsg is the error message displayed when a user is not found
	UserNotFoundErrMsg = "failed to get a Firebase user"

	// CustomTokenErrMsg is the error message displayed when a
	// custom firebase token is not created
	CustomTokenErrMsg = "failed to create custom token"

	// AuthenticateTokenErrMsg is the error message displayed when a
	// custom firebase token is not authenticated
	AuthenticateTokenErrMsg = "failed to authenticate custom token"

	// UpdateProfileErrMsg is the error message displayed when a
	// user profile is not found
	UpdateProfileErrMsg = "failed to update a user profile"

	// AddRecordErrMsg is the error message displayed when a
	// record fails to be added to the database
	AddRecordErrMsg = "failed to add the record to the database"

	// LikelyToRecommendErrMsg is the error message displayed that
	// occurs when the recommendation threshold is crossed
	LikelyToRecommendErrMsg = "the likelihood of recommending should be an int between 0 and 10"
)
