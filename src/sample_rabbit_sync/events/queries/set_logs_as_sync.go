package queries

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetLogsAsSync(logCollection *mongo.Collection, lastUnprocessedLogs []LastUnprocessedLog, currentTms int64) {
	log.Print("\n in UpdateEventTable func")
}
