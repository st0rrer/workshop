package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/st0rrer/workshop/web-service/src/config"
	"log"
	"net/http"
	"os"
)

type report interface {
	prepareMessages(stats []map[string]interface{}) ([]kafka.Message, error)
	sendMessages(ctx context.Context, messages []kafka.Message) error
}

var reports = make(map[string]report)

func GetReportHandler() func(res http.ResponseWriter, req *http.Request) {

	logger := log.New(os.Stderr, "GetReportHandler ", log.LstdFlags)

	activity, err := newActivity(config.GetConfig().ActivityTopic)
	if err != nil {
		logger.Println("Failed to create new activity", err)
	} else {
		reports["activity"] = activity
	}

	visit, err := newVisit(config.GetConfig().VisitTopic)
	if err != nil {
		logger.Println("Failed to create new visit", err)
	} else {
		reports["visit"] = visit
	}

	return reportHandler
}

func reportHandler(res http.ResponseWriter, req *http.Request) {

	logger := log.New(os.Stderr, "reportHandler ", log.LstdFlags)

	logger.Println(req.RequestURI)

	vars := mux.Vars(req)
	reportKey := vars["report"]

	var stats []map[string]interface{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&stats)
	if err != nil {
		logger.Println("error during deserialize ", err)
		http.Error(res, "invalid body request", 400)
	}

	report, ok := reports[reportKey]
	if !ok {
		http.Error(res, fmt.Sprintf("%v is not intialized", reportKey), 500)
		return
	}

	messages, err := report.prepareMessages(stats)
	if err != nil {
		logger.Println("error during prepare message ", err)
		http.Error(res, err.Error(), 500)
		return
	}

	err = report.sendMessages(req.Context(), messages)
	if err != nil {
		logger.Println("error during send message ", err)
		http.Error(res, err.Error(), 500)
		return
	}

	res.WriteHeader(200)
}
