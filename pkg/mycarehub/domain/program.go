package domain

// Program defines the program structure
type Program struct {
	ID           string       `json:"id"`
	Active       bool         `json:"active"`
	Name         string       `json:"name"`
	Organisation Organisation `json:"organisation"`
}

// ProgramPage returns a list of paginated programs
type ProgramPage struct {
	Pagination Pagination `json:"pagination"`
	Programs   []*Program `json:"programs"`
}
