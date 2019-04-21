package report

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/st0rrer/workshop/data-collector/src/config"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func Write(reportType string, messages []kafka.Message) error {

	data := []byte("[\n")
	for i, message := range messages {
		data = append(data, message.Value...)
		if i != len(messages)-1 {
			data = append(data, []byte(",")...)
			data = append(data, []byte("\n")...)
		}
	}
	data = append(data, []byte("\n]")...)

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	fileName := fmt.Sprintf("%v/%v-%v.json", config.GetConfig().RecordDirectory, strings.Title(reportType), timestamp)

	err := ioutil.WriteFile(fileName, data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
