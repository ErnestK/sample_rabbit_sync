package main

import (
	_ "fmt"
	"os"

	"github.com/streadway/amqp"
)

// {
//   "source": "snmp",
//   "component": "server1",
//   "resource": "CPU",
//   "crit": 3,
//   "message": "CPU Load > 80%",
//   "timestamp": 123456789
// }

// type logMessage struct {
// 	source    string
// 	component string
// 	crit      int
// 	message   string
// 	timestamp string
// }

//  !!! Very dirty file, need only for test sending message

func main() {
	// Get the connection string from the environment variable
	url := os.Getenv("AMQP_URL")

	//If it doesnt exist, use the default connection string
	if url == "" {
		url = "amqp://test:test@localhost:5672"
	}

	// Connect to the rabbitMQ instance
	connection, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}

	// Create a channel from the connection. We'll use channels to access the data in the queue rather than the
	// connection itself
	channel, err := connection.Channel()

	if err != nil {
		panic("could not open RabbitMQ channel:" + err.Error())
	}

	// We create an exahange that will bind to the queue to send and receive messages
	err = channel.ExchangeDeclare("canopsis.events", "topic", true, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	// We create a message to be sent to the queue.
	// It has to be an instance of the aqmp publishing struct

	message := amqp.Publishing{
		Body: []byte("Hello World"),
	}

	// We publish the message to the exahange we created earlier
	err = channel.Publish("canopsis.events", "random-key", false, false, message)

	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	// We create a queue named Test
	_, err = channel.QueueDeclare("Alarm", true, false, false, false, nil)

	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	err = channel.QueueBind("Alarm", "#", "canopsis.events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}
}
