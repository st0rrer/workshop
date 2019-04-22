package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/st0rrer/workshop/data-collector/src/config"
	"github.com/st0rrer/workshop/data-collector/src/mq"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

type MockConsumer struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	messages   chan []kafka.Message
	once       sync.Once
}

func (m *MockConsumer) BufReadMessages(ctx context.Context, buf int, timeout time.Duration) <-chan []kafka.Message {

	m.once.Do(func() {
		go func() {
			m.messages <- []kafka.Message{kafka.Message{Value: []byte("Test Message")}}
			close(m.messages)
			m.cancelFunc()
		}()
	})

	return m.messages
}

func (*MockConsumer) ReadMessage(ctx context.Context) <-chan kafka.Message {
	panic("implement me")
}

func (*MockConsumer) Commit(ctx context.Context, messages []kafka.Message) {
	log.Println("Commit consumer")
}

func (*MockConsumer) Close() {
	log.Println("Closing consumer")
}

func TestReport_ProcessMessages(t *testing.T) {

	dir := config.GetConfig().RecordDirectory
	err := prepareDirectory(dir)
	if err != nil {
		t.Fatalf("failed to prepare directory")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())

	report := Report{
		GroupID:  "groupId",
		Topic:    "topic",
		ctx:      ctx,
		consumer: mockConsumer(ctx, cancelFunc),
	}

	go func() {
		select {
		case <-time.After(30 * time.Second):
			t.Fatalf("Failed to process message")
		}
	}()

	report.ProcessMessages(1)

	files, _ := ioutil.ReadDir(dir)
	if len(files) != 1 {
		t.Errorf("Failed to save message")
	}
}

func mockConsumer(ctx context.Context, cancelFunc context.CancelFunc) mq.Consumer {
	return &MockConsumer{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		messages:   make(chan []kafka.Message),
	}
}

func prepareDirectory(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
