package queries

// Delimiter for CompositeGroupKey.
// In CompositeGroupKey I strconv two fields Component and Resource.
// In that way we can group by 2 fields but use only one( below I described why I do this)
var Delimiter = "-!-"

// EventLog data from test_technique_log.
// I add CompositeGroupKey cause grouping by multiple fields its normal but get data from that grouping this is headeachs. Maybe fix later
// Synchronized uses for set log message sync after we write it in test_technique
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
