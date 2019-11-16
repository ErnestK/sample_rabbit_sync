package main

import (
	"context"
	"log"
	"os"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	RabbitMQConnection *amqp.Connection
	EventChannel       <-chan amqp.Delivery
	MongoDBConnection  *mongo.Client
}

func GetConfig() *Config {
	connection, eventChan := getRabbitMQConfigs()
	return &Config{connection, eventChan, getMongoDBConnection()}
}

func getRabbitMQConfigs() (*amqp.Connection, <-chan amqp.Delivery) {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	//If it doesnt exist, use the default connection string
	if url == "" {
		url = "amqp://test:test@localhost:5672"
	}

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)
	// defer connection.Close()

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	log.Print("Connected to RabbitMQ!")

	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	_, err = channel.QueueDeclare(
		"canopsis.events", // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		panic("Failed to declare a queue:" + err.Error())
	}

	msgs, err := channel.Consume(
		"Alarm", // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		panic("error consuming the queue: " + err.Error())
	}

	return connection, msgs
}

func getMongoDBConnection() *mongo.Client {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	log.Print("Connected to MongoDB!")

	return client
}

func (config *Config) closeAllConnection() {
	config.RabbitMQConnection.Close()
	log.Print("Connection to RabbitMQ closed.")

	err := config.MongoDBConnection.Disconnect(context.TODO())

	if err != nil {
		panic(err)
	}
	log.Print("Connection to MongoDB closed.")
}
