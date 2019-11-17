package queries

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetEventIfExists return event from db if exists or false in bool return
func GetEventIfExists(component string, resource string, eventCollection *mongo.Collection) (bool, eventDb) {
	var fullEventMessage bson.M
	filter := bson.M{"component": component, "resource": resource}
	err := eventCollection.FindOne(context.TODO(), filter).Decode(&fullEventMessage)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, eventDb{}
		}
		log.Fatal(err)
	}
	return true, eventDb{ID: fullEventMessage["_id"].(primitive.ObjectID), Crit: fullEventMessage["crit"].(int32), Status: fullEventMessage["status"].(string)}
}
