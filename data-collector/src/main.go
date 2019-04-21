package main

import (
	"context"
	"github.com/st0rrer/workshop/data-collector/src/config"
	"log"
	"os"
	"os/signal"
	"sync"
)

func main() {

	logger := log.New(os.Stderr, "main ", log.LstdFlags)

	dir := config.GetConfig().RecordDirectory
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			logger.Panicf("failed to create directory %v for saving reports %v", dir, err)
		}
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for _, group := range config.GetConfig().Groups {
		wg.Add(1)

		go func(group config.Group) {
			logger.Println("start processing group:", group)
			defer wg.Done()

			report, err := NewReport(group.ID, group.Topic, ctx)
			if err != nil {
				logger.Println("failed to create report for topic: ", err)
				return
			}

			report.ProcessMessages(config.GetConfig().ConcurrentCount)
		}(group)
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
