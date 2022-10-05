package comms

// API response status
type Status string

const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusError   Status = "error"
)

// IsValid returns true if a status is valid
func (s Status) IsValid() bool {
	switch s {
	case "":
		return true
	}
	return false
}

// String representation of status
func (s Status) String() string {
	return string(s)
}
