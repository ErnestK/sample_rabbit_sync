package events

type EventStatus string

const (
	ONGOING  EventStatus = "ONGOING"
	RESOLVED EventStatus = "RESOLVED"
)

type event struct {
	component  string
	resource   string
	crit       int32
	last_msg   string
	first_msg  string
	start_time int32
	last_time  int32
	status     EventStatus
}
