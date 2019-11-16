package main

import (
	"log"
	"os"
	"os/signal"
	"sample_rabbit_sync/events"
	"sync"
	"syscall"
	"time"
)

// we can configure it and set to 0 if dont want any delay( but we can ddos db in that case)
var READ_TIMEOUT = 500 * time.Millisecond

func main() {
	var wgAddToBuffer sync.WaitGroup
	var wgSyncBuffer sync.WaitGroup

	config := GetConfig()

	ticker := time.NewTicker(READ_TIMEOUT)
	done := make(chan bool)
	exit := make(chan bool)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// We execute it in goroutine cause we want to use ticker because we dont want to ddos db
	go func() {
		wgSyncBuffer.Add(1)
		go events.SyncBuffer(config.MongoDBConnection, &wgSyncBuffer, done)

		for msg := range config.EventChannel {
			select {
			case <-done:
				return
			case <-ticker.C:
				wgAddToBuffer.Add(1)
				go events.AddToBuffer(config.MongoDBConnection, msg, &wgAddToBuffer)
			}
		}
	}()

	// Here we catch signals, finish all works and preapre to exit
	// we have 10 sec before system call 'kill' signal and exit prodram anywhere, all goroutines should done all work
	go func() {
		sig := <-sigs
		log.Print(sig)
		ticker.Stop()
		done <- true
		wgAddToBuffer.Wait()
		wgSyncBuffer.Wait()
		config.closeAllConnection()

		exit <- true
	}()

	<-exit
}
