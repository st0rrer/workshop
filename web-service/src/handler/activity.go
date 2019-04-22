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

type activity struct {
	topic    string
	producer mq.Producer
}

func newActivity(topic string) (*activity, error) {

	err := mq.CreateTopics(config.GetConfig().Brokers[0], topic)
	if err != nil {
		return nil, errors.Wrap(err, "can not create topics: ")
	}

	producer, err := mq.NewProducer(config.GetConfig().Brokers, topic)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create producer")
	}
	return &activity{
		topic:    topic,
		producer: producer,
	}, nil

}

var activityRules = govalidator.MapData{
	"dataVer":        []string{"required"},
	"userId":         []string{"required"},
	"activityType":   []string{"required"},
	"startTime":      []string{"required"},
	"endTime":        []string{"required"},
	"startLatitude":  []string{"required", "float"},
	"startLongitude": []string{"required", "float"},
	"endLatitude":    []string{"required", "float"},
	"endLongitude":   []string{"required", "float"},
}

func (a *activity) prepareMessages(stats []map[string]interface{}) ([]kafka.Message, error) {

	var messages []kafka.Message

	for _, s := range stats {
		err := a.validate(&s)
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

func (a *activity) validate(activity *map[string]interface{}) error {

	opts := govalidator.Options{
		Data:  activity,
		Rules: activityRules,
	}

	validator := govalidator.New(opts)
	e := validator.ValidateStruct()
	if len(e) > 0 {
		return errors.New(e.Encode())
	}

	return nil
}

func (a *activity) sendMessages(ctx context.Context, messages []kafka.Message) error {

	err := a.producer.PushMessage(ctx, messages)
	if err != nil {
		return errors.Wrap(err, "")
	}

	return nil

}
