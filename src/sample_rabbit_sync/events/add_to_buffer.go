package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"sync"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

func AddToBuffer(mongoClient *mongo.Client, msg amqp.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Print("message received: " + string(msg.Body))
	event, err := deserializeMessageFromMQ(msg.Body)

	if err != nil {
		log.Fatal("error during deserialize message: " + string(msg.Body))
	}

	log.Print("deserilized message source: " + event.Source)

	collection := mongoClient.Database("test").Collection("test_technique_log")
	insertResult, err := collection.InsertOne(context.TODO(), event)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Print("Inserted a single document: ", insertResult.InsertedID)
	}

	// TODO: select all for debug
	cursor, err := collection.Find(context.TODO(), bson.D{})
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
				os.Exit(1)

				// If there are no cursor.Decode errors
			} else {
				fmt.Println("\nresult type:", reflect.TypeOf(result))
				fmt.Println("result:", result)
			}
		}
	}

	msg.Ack(true)
}

func deserializeMessageFromMQ(b []byte) (event, error) {
	var msg event
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}
