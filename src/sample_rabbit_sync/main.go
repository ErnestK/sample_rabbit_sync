package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/streadway/amqp"
	_ "go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

var READ_TIMEOUT = 500 * time.Millisecond

type mqMessage struct {
	Source    string
	Component string
	Crit      int
	Message   string
	Timestamp int32
}

func main() {
	config := GetConfig()

	ticker := time.NewTicker(READ_TIMEOUT)
	done := make(chan bool)
	exit := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// We execute it in goroutine cause we want to use ticker because we dont want to ddos db
	go func() {
		for msg := range config.EventChannel {
			select {
			case <-done:
				return
			case <-ticker.C:

				log.Print("message received: " + string(msg.Body))
				logMsg, err := deserializeMessageFromMQ(msg.Body)

				if err != nil {
					log.Fatal("error during deserialize message: " + string(msg.Body))
				}

				log.Print("deserilized message source: " + logMsg.Source)
				msg.Ack(true)
			}
		}
	}()

	// collection := config.MongoDBConnection.Database("test").Collection("test_technique")
	// insertResult, err := collection.InsertOne(context.TODO(), logMsg)
	// if err != nil {
	// 	// We log error, and retry later
	// 	log.Fatal(err)
	// } else {
	// 	log.Print("Inserted a single document: ", insertResult.InsertedID)
	// 	msg.Ack(true)
	// }

	// Here we catch signals finish all works and preapre to exit
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		closeAllConnection(config)
		exit <- true
		done <- true
		log.Print("wirte exit true")
	}()

	log.Print("before read from exit")
	<-exit
	log.Print("after read from exit")
}

func deserializeMessageFromMQ(b []byte) (mqMessage, error) {
	var msg mqMessage
	buf := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(&msg)
	return msg, err
}

func closeAllConnection(config *Config) {
	config.RabbitMQConnection.Close()
	log.Print("Connection to RabbitMQ closed.")

	err := config.MongoDBConnection.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}
	log.Print("Connection to MongoDB closed.")
}
