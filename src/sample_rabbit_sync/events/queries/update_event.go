package queries

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateEventTable updated data
func UpdateEventTable(eventCollection *mongo.Collection, row eventDb, fullLogMessage primitive.M) {
	log.Print("\n in UpdateEventTable func")
	newCrit := fullLogMessage["crit"].(int32)
	newMessage := fullLogMessage["message"].(string)
	newTimestmap := fullLogMessage["timestamp"].(int32)

	filter := bson.D{{"_id", row.ID}}

	// TODO: here many updates because I tired fight with mongo( maybe nobody come here and dont see that shame)
	update1 := bson.D{{"$set", bson.D{{"crit", newCrit}}}}
	updateOne(eventCollection, filter, update1)
	update2 := bson.D{{"$set", bson.D{{"last_msg", newMessage}}}}
	updateOne(eventCollection, filter, update2)
	update3 := bson.D{{"$set", bson.D{{"last_time", newTimestmap}}}}
	updateResult := updateOne(eventCollection, filter, update3)

	if newCrit == 0 {
		update4 := bson.D{{"$set", bson.D{{"status", RESOLVED}}}}
		updateOne(eventCollection, filter, update4)
	}
	log.Print("\nUpdated documents ", updateResult.ModifiedCount)
}

func updateOne(eventCollection *mongo.Collection, filter primitive.D, update primitive.D) *mongo.UpdateResult {
	updateResult, err := eventCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	return updateResult
}
