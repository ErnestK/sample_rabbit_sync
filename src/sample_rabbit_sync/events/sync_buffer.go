package events

import (
	"log"
	"sync"
	"time"

	_ "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	_ "go.mongodb.org/mongo-driver/mongo/options"
)

var SYNC_TIMEOUT = 500 * time.Millisecond

func SyncBuffer(mongoClient *mongo.Client, wg *sync.WaitGroup, done <-chan bool) {
	ticker := time.NewTicker(SYNC_TIMEOUT)

	for {
		select {
		case <-done:
			ticker.Stop()
			wg.Done()
			return
		case t := <-ticker.C:
			log.Print("Tick at sync", t)
		}
	}
}
