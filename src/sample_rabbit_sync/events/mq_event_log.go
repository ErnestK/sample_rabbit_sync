package events

type mqEventLog struct {
	Source    string
	Component string
	Resource  string
	Crit      int
	Message   string
	Timestamp int32
}
