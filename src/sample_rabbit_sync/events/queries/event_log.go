package queries

var Delimiter = "-!-"

type EventLog struct {
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
