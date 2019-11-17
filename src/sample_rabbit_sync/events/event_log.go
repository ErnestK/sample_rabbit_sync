package events

var delimiter = "-!-"

type mqEventLog struct {
	Source    string
	Component string
	Resource  string
	Crit      int
	Message   string
	Timestamp int32
}

type eventLog struct {
	Source            string
	Component         string
	Resource          string
	Crit              int
	Message           string
	Timestamp         int32
	CompositeGroupKey string
	Synchronized      bool
	CreatedAt         int64
}
