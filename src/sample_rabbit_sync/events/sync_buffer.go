package events

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sample_rabbit_sync/events/queries"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
		case t := <-ticker.C:
			log.Print("Tick at sync", t)

			logRows := queries.GetLastUnprocessedRows(logCollection, currentTms)
			syncWithMainTable(mongoClient, logRows, currentTms)
		}
	}
}

func syncWithMainTable(mongoClient *mongo.Client, lastUnprocessedLogs []queries.LastUnprocessedLog, currentTms int64) {
	eventCollection := mongoClient.Database("test").Collection("test_technique")
	logCollection := mongoClient.Database("test").Collection("test_technique_log")

	// TODO: select all for debug
	cursor, err := logCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		fmt.Println("Finding all documents ERROR:", err)
		defer cursor.Close(context.TODO())

		// If the API call was a success
	} else {
		for cursor.Next(context.TODO()) {
			var result bson.M
			err := cursor.Decode(&result)

			// If there is a cursor.Decode error
			if err != nil {
				fmt.Println("cursor.Next() error:", err)

				// If there are no cursor.Decode errors
			} else {
				fmt.Println("\nresult type:", reflect.TypeOf(result))
				fmt.Println("result:", result)
			}
		}
	}

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
			event := event{
				component:  fullLogMessage["component"].(string),
				resource:   fullLogMessage["resource"].(string),
				crit:       fullLogMessage["crit"].(int32),
				last_msg:   fullLogMessage["message"].(string),
				first_msg:  fullLogMessage["message"].(string),
				start_time: fullLogMessage["timestamp"].(int32),
				last_time:  fullLogMessage["timestamp"].(int32),
				status:     ONGOING,
			}

			insertResult, err := eventCollection.InsertOne(context.TODO(), event)

			if err != nil {
				log.Fatal(err)
			} else {
				log.Print("\nInserted a single document in test_technique table: ", insertResult.InsertedID)
			}
		}
		queries.SetLogsAsSync(logCollection, lastUnprocessedLogs, currentTms)
	}
}
