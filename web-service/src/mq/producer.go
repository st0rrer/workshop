package mq

import (
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func init() {
	flag.IntVar(&numPartitions, "num-partitions", 4, "default partitions for new topic")
	flag.IntVar(&replicationFactor, "replication-factor", 2, "default replication factor for new topic")
	flag.Parse()
}

var numPartitions int
var replicationFactor int

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) (producer *Producer, err error) {

	defer func() {
		if r := recover(); r != nil {
			producer, err = nil, errors.Errorf("%v", r)
		}
	}()

	cfg := kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	}

	producer, err = &Producer{
		writer: kafka.NewWriter(cfg),
	}, nil

	return
}

func (p *Producer) PushMessage(ctx context.Context, msg []kafka.Message) error {
	return p.writer.WriteMessages(ctx, msg...)
}

func CreateTopics(broker string, topic string) error {
	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		return errors.Wrapf(err, "could not connect to %v", broker)
	}
	defer conn.Close()

	controllerBroker, err := conn.Controller()
	if err != nil {
		return errors.Wrap(err, "could not find controllerBroker")
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%v:%v", controllerBroker.Host, controllerBroker.Port))
	if err != nil {
		return errors.Wrapf(err, "could not connect to controllerBroker %v", controllerBroker)
	}
	defer controllerConn.Close()

	topicConfigs := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}

	err = controllerConn.CreateTopics(topicConfigs)
	if err != nil {
		return errors.Wrapf(err, "could not create topic %v", topic)
	}

	return nil
}
