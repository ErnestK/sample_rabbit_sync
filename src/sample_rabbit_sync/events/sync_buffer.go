package events

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var SYNC_TIMEOUT = 1000 * time.Millisecond

func SyncBuffer(mongoClient *mongo.Client, wg *sync.WaitGroup, done <-chan bool) {
	ticker := time.NewTicker(SYNC_TIMEOUT)
	// eventCollection := mongoClient.Database("test").Collection("test_technique")
	logCollection := mongoClient.Database("test").Collection("test_technique_log")

	for {
		// currentTms := time.Now().UnixNano()

		select {
		case <-done:
			ticker.Stop()
			wg.Done()
			return
		case t := <-ticker.C:
			log.Print("Tick at sync", t)
			// start

			// specify a pipeline that will return the number of times each name appears in the collection
			// specify the MaxTime option to limit the amount of time the operation can run on the server
			groupBy := bson.D{
				{"$group", bson.D{
					{"_id", bson.E{"$mqeventlog.component", "$mqeventlog.resource"}},
					{"maxTimestamp", bson.D{
						{"$max", "$mqeventlog.timestamp"},
					}},
				}},
			}
			opts := options.Aggregate().SetMaxTime(2 * time.Second)
			cursor, err := logCollection.Aggregate(context.TODO(), mongo.Pipeline{groupBy}, opts)
			if err != nil {
				log.Fatal(err)
			}

			// get a list of all returned documents and print them out
			// see the mongo.Cursor documentation for more examples of using cursors
			var results []bson.M
			if err = cursor.All(context.TODO(), &results); err != nil {
				log.Fatal(err)
			}
			for _, result := range results {
				fmt.Printf("name %v appears %v times\n", result["_id"], result["maxTimestamp"])
			}
			// end
		}
	}
}

// If it has the attribute crit > 0 :
// – If an alarm already exists in the database for the same component / resource couple
// with status = ONGOING :
// *
// Then update crit, last_msg and last_time fields.
// – If no alarm is present in the database for the same component / resource couple with
// status = ONGOING :
// *
// Then insert an alarm in the database with the following attributes :
//
// Alarm object fields Event object fields
// component component
//
//
// Alarm object fields Event object fields
// resource resource
// crit crit
// last_msg message
// first_msg message
// start_time timestamp
// last_time timestamp
// status ONGOING (fixed)
//
//
//
// • If the event has the attribute crit = 0 :
// – If an alarm already exists in the database for the same component / resource couple
// with status = ONGOING :
// *
// Then update crit, last_msg and last_time fields and change status to
// RESOLVED.
// – If no alarm is present in the database for the same component / resource couple with
// status = ONGOING :
// *
// Go to next message.
