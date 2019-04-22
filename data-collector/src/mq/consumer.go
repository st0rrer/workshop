package mq

import (
	"context"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Consumer interface {
	BufReadMessages(ctx context.Context, buf int, timeout time.Duration) <-chan []kafka.Message
	ReadMessage(ctx context.Context) <-chan kafka.Message
	Commit(ctx context.Context, messages []kafka.Message)
	Close()
}

type KafkaConsumer struct {
	reader      *kafka.Reader
	readOnce    sync.Once
	bufReadOnce sync.Once
	message     <-chan kafka.Message
	bufMessages <-chan []kafka.Message
}

func NewConsumer(brokers []string, topic string, groupID string) (consumer Consumer, err error) {

	defer func() {
		if r := recover(); r != nil {
			consumer, err = nil, errors.Errorf("%v", r)
		}
	}()

	cfg := kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	}

	consumer = &KafkaConsumer{
		reader: kafka.NewReader(cfg),
	}
	return
}

func (c *KafkaConsumer) BufReadMessages(ctx context.Context, buf int, timeout time.Duration) <-chan []kafka.Message {

	c.bufReadOnce.Do(func() {
		c.bufMessages = c.bufReadMessages(ctx, buf, timeout)
	})

	return c.bufMessages
}

func (c *KafkaConsumer) bufReadMessages(ctx context.Context, buf int, timeout time.Duration) <-chan []kafka.Message {
	messages := make(chan []kafka.Message)
	go func() {

		defer func() {
			close(messages)
		}()

		var bufMessages []kafka.Message
		timer := time.NewTimer(timeout)
		for {
			select {
			case msg, ok := <-c.ReadMessage(ctx):

				if !ok {
					return
				}

				bufMessages = append(bufMessages, msg)

				if len(bufMessages) == buf {
					messages <- bufMessages
					bufMessages = nil
					timer.Reset(timeout)
				}

			case <-ctx.Done():

				if len(messages) > 0 {
					messages <- bufMessages
				}

				return
			case <-timer.C:

				if len(bufMessages) > 0 {
					messages <- bufMessages
					bufMessages = nil
				}

				timer.Reset(timeout)
			}
		}
	}()

	return messages
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) <-chan kafka.Message {

	c.readOnce.Do(func() {
		c.message = c.readMessage(ctx)
	})

	return c.message

}

func (c *KafkaConsumer) readMessage(ctx context.Context) <-chan kafka.Message {
	logger := log.New(os.Stderr, "consumer#readMessage ", log.LstdFlags)

	messages := make(chan kafka.Message)
	go func() {
		defer close(messages)

		for {
			message, err := c.reader.FetchMessage(ctx)
			if err != nil {
				switch err {
				case io.EOF:
					logger.Println("reader is closed")
					return
				case ctx.Err():
					logger.Println("context in closed")
					c.Close()
					return
				default:
					logger.Println("error during fetchMessage", err)
					continue
				}
			}

			messages <- message
		}
	}()

	return messages
}

func (c *KafkaConsumer) Commit(ctx context.Context, messages []kafka.Message) {
	logger := log.New(os.Stderr, "consumer#Commit ", log.LstdFlags)
	err := c.reader.CommitMessages(ctx, messages...)
	if err != nil {
		logger.Println("error during committing messages err", err)
	}
}

func (c *KafkaConsumer) Close() {
	logger := log.New(os.Stderr, "consumer#Close ", log.LstdFlags)
	err := c.reader.Close()
	if err != nil {
		logger.Println("error during closing reader", err)
	}
}
