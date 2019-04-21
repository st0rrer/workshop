package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/st0rrer/workshop/data-collector/src/config"
	"github.com/st0rrer/workshop/data-collector/src/mq"
	"github.com/st0rrer/workshop/data-collector/src/report"
	"log"
	"os"
	"sync"
	"time"
)

type Report struct {
	GroupID    string
	Topic      string
	writeError chan error
	once       sync.Once
	wg         sync.WaitGroup
	ctx        context.Context
	consumer   *mq.Consumer
}

func NewReport(groupID string, topic string, ctx context.Context) (*Report, error) {

	messageQueue, err := mq.NewConsumer(config.GetConfig().Brokers, topic, groupID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new consumer: %v", err)
	}

	return &Report{
		GroupID:    groupID,
		Topic:      topic,
		ctx:        ctx,
		consumer:   messageQueue,
		writeError: make(chan error, 100),
	}, nil
}

func (r *Report) ProcessMessages(concurrentCount int) {

	for i := 0; i < concurrentCount; i++ {
		r.wg.Add(1)
		go r.saveRecord()
	}

	r.wg.Wait()
	r.consumer.Close()
}

func (r *Report) saveRecord() {

	logger := log.New(os.Stderr, "report#save ", log.LstdFlags)

	defer r.wg.Done()

	for {
		select {
		case messages, ok := <-r.consumer.BufReadMessages(r.ctx, 100, 2*time.Minute):

			if len(messages) == 0 && !ok {
				return
			}

			err := report.Write(r.Topic, messages)
			if err != nil {
				logger.Println("failed writing to file", err)
			} else {
				r.consumer.Commit(r.ctx, messages)
			}
		}
	}
}
