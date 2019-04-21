package config

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

func init() {

	flag.Var(&conf.Groups, "groups", "kafka groups")
	flag.Var(&conf.Brokers, "brokers", "kafka brokers")
	flag.StringVar(&conf.RecordDirectory, "record-directory", "./log/", "path to save record")
	flag.IntVar(&conf.ConcurrentCount, "concurrent-count", 10, "count of goroutines for processing messages")

	flag.Parse()
}

var conf = Config{}

type Config struct {
	Groups          Groups
	Brokers         Brokers
	RecordDirectory string
	ConcurrentCount int
}

func GetConfig() Config {
	return conf
}

type Group struct {
	ID    string
	Topic string
}

func (g *Group) String() string {
	return fmt.Sprintf("group name: %v, topic name: %v", g.ID, g.Topic)
}

type Groups []Group

func (t *Groups) String() string {
	groups := []Group(*t)
	var res []string
	for _, group := range groups {
		res = append(res, group.String())
	}
	return strings.Join(res, ", ")
}

func (t *Groups) Set(value string) error {
	matched, err := regexp.MatchString("([a-zA-Z0-9\\._\\-]*):([a-zA-Z0-9\\._\\-]*)", value)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("Incorrect group")
	}
	groupValues := strings.Split(value, ":")
	*t = append(*t, Group{
		ID:    groupValues[0],
		Topic: groupValues[1],
	})
	return nil
}

type Brokers []string

func (b *Brokers) String() string {
	return strings.Join([]string(*b), ",")
}

func (b *Brokers) Set(value string) error {
	*b = append(*b, value)
	return nil
}
