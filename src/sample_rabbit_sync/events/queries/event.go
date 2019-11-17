package queries

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type eventStatus string

// its enum for event status
const (
	ONGOING  eventStatus = "ONGOING"
	RESOLVED eventStatus = "RESOLVED"
)

// Event data from test_technique
type Event struct {
	Component  string
	Resource   string
	Crit       int32
	Last_msg   string
	First_msg  string
	Start_time int32
	Last_time  int32
	Status     eventStatus
}

type eventDb struct {
	ID     primitive.ObjectID
	Crit   int32
	Status string
}
