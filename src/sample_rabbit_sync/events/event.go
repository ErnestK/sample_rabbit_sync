package events

type event struct {
	Source    string
	Component string
	Resource  string
	Crit      int
	Message   string
	Timestamp int32
}
