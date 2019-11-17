package events

import (
	"context"
	"log"
	"sample_rabbit_sync/events/queries"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var SYNC_TIMEOUT = 1000 * time.Millisecond

func SyncBuffer(mongoClient *mongo.Client, wg *sync.WaitGroup, done <-chan bool) {
	ticker := time.NewTicker(SYNC_TIMEOUT)
	logCollection := mongoClient.Database("test").Collection("test_technique_log")

	for {
		currentTms := time.Now().UnixNano()

		select {
		case <-done:
			ticker.Stop()
			wg.Done()
			return
		case <-ticker.C:
			logRows := queries.GetLastUnprocessedRows(logCollection, currentTms)
			syncWithMainTable(mongoClient, logRows, currentTms)
		}
	}
}

func syncWithMainTable(mongoClient *mongo.Client, lastUnprocessedLogs []queries.LastUnprocessedLog, currentTms int64) {
	eventCollection := mongoClient.Database("test").Collection("test_technique")
	logCollection := mongoClient.Database("test").Collection("test_technique_log")

	for _, lastUnprocessedLog := range lastUnprocessedLogs {
		var fullLogMessage bson.M

		logFilter := bson.M{"component": lastUnprocessedLog.Component, "resource": lastUnprocessedLog.Resource, "timestamp": lastUnprocessedLog.MaxTimestamp}
		err := logCollection.FindOne(context.TODO(), logFilter).Decode(&fullLogMessage)
		if err != nil {
			log.Fatal(err)
		}

		isExists, row := queries.GetEventIfExists(lastUnprocessedLog.Component, lastUnprocessedLog.Resource, eventCollection)
		if isExists {
			queries.UpdateEventTable(eventCollection, row, fullLogMessage)
		} else {
			event := queries.Event{
				Component:  fullLogMessage["component"].(string),
				Resource:   fullLogMessage["resource"].(string),
				Crit:       fullLogMessage["crit"].(int32),
				Last_msg:   fullLogMessage["message"].(string),
				First_msg:  fullLogMessage["message"].(string),
				Start_time: fullLogMessage["timestamp"].(int32),
				Last_time:  fullLogMessage["timestamp"].(int32),
				Status:     queries.ONGOING,
			}

			insertResult, err := eventCollection.InsertOne(context.TODO(), event)

			if err != nil {
				log.Fatal(err)
			} else {
				log.Print("\nInserted a single document in test_technique table: ", insertResult.InsertedID)
			}
		}
		queries.SetLogAsSync(logCollection, fullLogMessage["_id"].(primitive.ObjectID))
	}
}
