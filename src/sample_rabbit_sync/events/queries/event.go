package queries

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventStatus string

const (
	ONGOING  EventStatus = "ONGOING"
	RESOLVED EventStatus = "RESOLVED"
)

type Event struct {
	Component  string
	Resource   string
	Crit       int32
	Last_msg   string
	First_msg  string
	Start_time int32
	Last_time  int32
	Status     EventStatus
}

type eventDb struct {
	ID     primitive.ObjectID
	Crit   int32
	Status string
}
