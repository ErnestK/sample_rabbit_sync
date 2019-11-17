package queries

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEventIfExists(component string, resource string, eventCollection *mongo.Collection) (bool, eventDb) {
	var fullEventMessage bson.M
	filter := bson.M{"component": component, "resource": resource}
	err := eventCollection.FindOne(context.TODO(), filter).Decode(&fullEventMessage)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, eventDb{}
		} else {
			log.Fatal(err)
		}
	}
	return true, eventDb{Crit: fullEventMessage["crit"].(int), Status: fullEventMessage["status"].(string)}
}

type eventDb struct {
	Crit   int
	Status string
}