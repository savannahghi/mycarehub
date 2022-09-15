package dto

// WorkStationDetailsOutput contains the details of specific items in a facility.
//
// These include things like number of notification associated to that facility, client's surveys, service requests etc
type WorkStationDetailsOutput struct {
	Notifications   int `json:"notifications"`
	Surveys         int `json:"surveys"`
	Articles        int `json:"articles"`
	Messages        int `json:"messages"`
	ServiceRequests int `json:"serviceRequests"`
}
