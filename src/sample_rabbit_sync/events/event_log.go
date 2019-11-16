package events

type mqEventLog struct {
	Source    string
	Component string
	Resource  string
	Crit      int
	Message   string
	Timestamp int32
}

type eventLog struct {
	MqEventLog   mqEventLog
	Synchronized bool
	CreatedAt    int64
}
