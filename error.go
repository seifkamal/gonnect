package gonnect

import "fmt"

type NoResultsFound struct{
	Query string
	Values map[string]interface{}
}

func (e NoResultsFound) Error() string {
	return fmt.Sprintf("no results found in storage. Query: '%s' Value %v", e.Query, e.Values)
}

type InsufficientPlayers struct {
	Need, Have int
}

func (e InsufficientPlayers) Error() string {
	return fmt.Sprintf("insufficient players to start a match. Need %v but only found %v", e.Need, e.Have)
}
