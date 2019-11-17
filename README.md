# sample_rabbit_send
dirty file, used as template need for sending message.

# sample_rabbit_sync

## Step to run 
1) Run rabbitMQ in container           
`docker run -d --hostname my-rabbit --name some-rabbit -e RABBITMQ_DEFAULT_USER=test -e RABBITMQ_DEFAULT_PASS=test  -p 5672:5672 -p 5673:5673 -p 15672:15672 rabbitmq:3-management`         
2) Run mongoDB in container         
`docker run -d -it -p 27017:27017 mongo`              

3) run main from bin `bin/main` or `go run main.go config.go`
