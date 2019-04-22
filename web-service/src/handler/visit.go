package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
	"github.com/st0rrer/workshop/web-service/src/config"
	"github.com/st0rrer/workshop/web-service/src/mq"
	"github.com/thedevsaddam/govalidator"
	"log"
)

type visit struct {
	topic    string
	producer mq.Producer
}

func newVisit(topic string) (*visit, error) {

	err := mq.CreateTopics(config.GetConfig().Brokers[0], topic)
	if err != nil {
		return nil, errors.Wrap(err, "can not create topics: ")
	}

	producer, err := mq.NewProducer(config.GetConfig().Brokers, topic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create producer")
	}
	return &visit{
		topic:    topic,
		producer: producer,
	}, nil

}

var visitRules = govalidator.MapData{
	"dataVer":       []string{"required"},
	"userId":        []string{"required"},
	"algorithmType": []string{"required"},
	"enterTime":     []string{"required"},
	"exitTime":      []string{"required"},
	"poiId":         []string{"required"},
	"latitude":      []string{"required", "float"},
	"longitude":     []string{"required", "float"},
}

func (v *visit) prepareMessages(stats []map[string]interface{}) ([]kafka.Message, error) {

	var messages []kafka.Message

	for _, s := range stats {
		err := v.validate(&s)
		if err != nil {
			log.Println("activity validation error", err)
			return nil, err
		}

		body, err := json.Marshal(s)
		if err != nil {
			log.Println("serializing error ", err)
			return nil, fmt.Errorf("problem, can not prepare message for mq")
		}

		messages = append(messages, kafka.Message{
			Value: body,
		})
	}

	return messages, nil
}

func (v *visit) validate(visit *map[string]interface{}) error {

	opts := govalidator.Options{
		Data:  visit,
		Rules: visitRules,
	}

	validator := govalidator.New(opts)
	e := validator.ValidateStruct()
	if len(e) > 0 {
		return errors.New(e.Encode())
	}

	return nil
}

func (v *visit) sendMessages(ctx context.Context, messages []kafka.Message) error {

	err := v.producer.PushMessage(ctx, messages)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}
