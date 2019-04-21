package config

import (
	"flag"
)

func init() {

	flag.IntVar(&conf.ConcurrentClients, "concurrent-clients", 2, "Concurrent Clients")
	flag.IntVar(&conf.ReportInterval, "report-interval", 10, "send report with interval in seconds")
	flag.StringVar(&conf.WebService, "web-service", "http://127.0.0.1:8080", "URL to send stats")

	flag.Parse()
}

var conf = Config{}

type Config struct {
	ConcurrentClients int
	ReportInterval    int
	WebService        string
}

func GetConfig() Config {
	return conf
}
