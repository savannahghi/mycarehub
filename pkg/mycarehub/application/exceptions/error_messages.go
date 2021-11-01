package exceptions

const (
	// NormalizeMSISDNErrMsg is the error message displayed when
	// normalize the msisdn(phone number) fails
	NormalizeMSISDNErrMsg = "unable to normalize the phonenumber"

	// PINMismatchErrMsg is the error message displayed when
	// the user supplied PIN number does not match the PIN
	// record we have stored
	PINMismatchErrMsg = "wrong PIN credentials supplied"

	// ExpiredPinErrMsg is the error message displayed when the user supplied pin
	// has expired
	ExpiredPinErrMsg = "the user pin has expired"
)
