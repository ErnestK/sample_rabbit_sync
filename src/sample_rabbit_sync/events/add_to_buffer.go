package events

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddToBuffer(mongoClient *mongo.Client, msg amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("message received: " + string(msg.Body))
	mqEventLog, err := deserializeMessageFromMQ(msg.Body)

	if err != nil {
		log.Fatal("error during deserialize message: " + string(msg.Body))
	}
	// using composite key because groupong by two column and after tring to decode data is hell
	eventLog := eventLog{
		Source:            mqEventLog.Source,
		Component:         mqEventLog.Component,
		Resource:          mqEventLog.Resource,
		Crit:              mqEventLog.Crit,
		Message:           mqEventLog.Message,
		Timestamp:         mqEventLog.Timestamp,
		Synchronized:      false,
		CreatedAt:         time.Now().UnixNano(),
		CompositeGroupKey: mqEventLog.Component + delimiter + mqEventLog.Resource,
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
