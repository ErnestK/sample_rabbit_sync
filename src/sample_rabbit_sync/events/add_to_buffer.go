package events

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sample_rabbit_sync/events/queries"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddToBuffer insert new data row from rabbit, into table test_technique_log.
// Since we insert only we dont have blocking/duplicates on db, and can use it in goroutines
// Later other functions handles data from that table and update/insert in final table
func AddToBuffer(mongoClient *mongo.Client, msg amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("message received: " + string(msg.Body))
	mqEventLog, err := deserializeMessageFromMQ(msg.Body)

	if err != nil {
		log.Fatal("error during deserialize message: " + string(msg.Body))
	}
	// using composite key because groupong by two column and after tring to decode data is hell
	eventLog := queries.EventLog{
		Source:            mqEventLog.Source,
		Component:         mqEventLog.Component,
		Resource:          mqEventLog.Resource,
		Crit:              mqEventLog.Crit,
		Message:           mqEventLog.Message,
		Timestamp:         mqEventLog.Timestamp,
		Synchronized:      false,
		CreatedAt:         time.Now().UnixNano(),
		CompositeGroupKey: mqEventLog.Component + queries.Delimiter + mqEventLog.Resource,
	}

	collection := mongoClient.Database("test").Collection("test_technique_log")
	insertResult, err := collection.InsertOne(context.TODO(), eventLog)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Inserted a single document: ", insertResult.InsertedID)
	}

	msg.Ack(true)
}

func deserializeMessageFromMQ(b []byte) (mqEventLog, error) {
	var msg mqEventLog
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
