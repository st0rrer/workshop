package config

import (
	"flag"
	"strings"
)

func init() {
	flag.Var(&conf.Brokers, "brokers", "kafka brokers")
	flag.IntVar(&conf.Port, "port", 8080, "http listening port")
	flag.StringVar(&conf.ActivityTopic, "activity-topic", "activity", "activity topic name")
	flag.StringVar(&conf.VisitTopic, "visit-topic", "visit", "visit topic name")
	flag.Parse()
}

var conf = Config{}

type Config struct {
	Port          int
	Brokers       Brokers
	ActivityTopic string
	VisitTopic    string
}

func GetConfig() Config {
	return conf
}

type Brokers []string

func (b *Brokers) String() string {
	return strings.Join([]string(*b), ",")
}

func (b *Brokers) Set(value string) error {
	*b = append(*b, value)
	return nil
}
