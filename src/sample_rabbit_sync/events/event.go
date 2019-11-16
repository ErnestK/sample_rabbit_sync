package events

type event struct {
	component  string
	resource   string
	crit       int
	last_msg   string
	first_msg  string
	start_time int
	last_time  int
	status     string
}
