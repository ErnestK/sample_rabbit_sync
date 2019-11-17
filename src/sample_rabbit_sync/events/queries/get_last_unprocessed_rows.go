package queries

import (
	"context"
	"log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// we can update and insert right in that method and this would be faster, but less readeable
func GetLastUnprocessedRows(logCollection *mongo.Collection, currentTms int64) []LastUnprocessedLog {
	delimiter := "-!-"
	// TODO: limit here maybe, depends of data
	pipe := mongo.Pipeline{
		{{"$match", bson.D{{"synchronized", false}}}},
		{{"$match", bson.D{{"createdat", bson.D{{"$lte", currentTms}}}}}},
		{{"$group", bson.D{
			{"_id", "$compositegroupkey"},
			{"maxTimestamp", bson.D{
				{"$max", "$timestamp"},
			}},
		}}},
	}

	opts := options.Aggregate().SetMaxTime(2 * time.Second)
	cursor, err := logCollection.Aggregate(context.TODO(), pipe, opts)
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	var lastUnprocessedLogsArr []LastUnprocessedLog
	for _, result := range results {
		if isPrimitiveNil(result["maxTimestamp"]) || isPrimitiveNil(result["_id"]) {
			log.Print("nil in result")
		} else {
			maxTimestamp := result["maxTimestamp"].(int32)
			compositeKeys := strings.Split(result["_id"].(string), delimiter)
			if len(compositeKeys) != 2 {
				log.Fatal("function to create composite key is broken, all bad!")
			}

			lastUnprocessedLogV := LastUnprocessedLog{Component: compositeKeys[0], Resource: compositeKeys[1], MaxTimestamp: maxTimestamp}

			lastUnprocessedLogsArr = append(lastUnprocessedLogsArr, lastUnprocessedLogV)
		}
	}

	return lastUnprocessedLogsArr
}

type LastUnprocessedLog struct {
	Component    string
	Resource     string
	MaxTimestamp int32
}

func isPrimitiveNil(pr interface{}) bool {
	return pr == nil || (reflect.ValueOf(pr).Kind() == reflect.Ptr && reflect.ValueOf(pr).IsNil())
}
