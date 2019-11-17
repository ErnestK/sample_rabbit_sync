package queries

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetLogAsSync(logCollection *mongo.Collection, logId primitive.ObjectID) {
	filter := bson.D{{"_id", logId}}
	update := bson.D{{"$set", bson.D{{"synchronized", true}}}}

	result, err := logCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	if result.MatchedCount != 0 {
		log.Print("matched and replaced an existing document: ", result.MatchedCount)
		return
	}
}
