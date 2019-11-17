package queries

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
