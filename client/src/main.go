package main

import (
	"context"
	"github.com/st0rrer/workshop/client/src/config"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {

	logger := log.New(os.Stderr, "main ", log.LstdFlags)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for i := 0; i < config.GetConfig().ConcurrentClients; i++ {

		wg.Add(1)
		go func() {
			defer wg.Done()

			client, err := NewClient(ctx)
			if err != nil {
				logger.Printf("can not create client %v", err)
				return
			}

			err = client.Start()
			if err != nil {
				logger.Printf("client error: %v", err)
			}
		}()
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancelFunc()
		wg.Wait()
		os.Exit(0)
	}()

	<-ctx.Done()
	logger.Println("Gracefully stop application")
}
