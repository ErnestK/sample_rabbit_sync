package queries

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
)

func SetLogsAsSync(logCollection *mongo.Collection, lastUnprocessedLogs []LastUnprocessedLog, currentTms int64) {
	log.Print("\n in UpdateEventTable func")
	// updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
}
