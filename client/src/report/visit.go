package report

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/st0rrer/workshop/client/src/config"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

type visit struct {
	DataVer       int       `json:"dataVer"`
	UserId        uuid.UUID `json:"userId"`
	AlgorithmType int       `json:"algorithmType"`
	EnterTime     timestamp `json:"enterTime"`
	ExitTime      timestamp `json:"exitTime"`
	PoiId         int64     `json:"poiId"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
}

func newVisit(userId uuid.UUID, startTime time.Time) visit {
	return visit{
		DataVer:       1,
		UserId:        userId,
		AlgorithmType: RandomInt(1, 8, 1),
		EnterTime:     timestamp(startTime),
		ExitTime:      timestamp(startTime.Add(40 * time.Minute)),
		PoiId:         rand.Int63(),
		Latitude:      RandomLatitude(),
		Longitude:     RandomLongitude(),
	}
}

type visitStats []visit

func (v visitStats) send() *statsResult {

	var logging = log.New(os.Stdout, "visit#send ", 0)

	res, err := doPost(v, fmt.Sprintf("%v/%v", config.GetConfig().WebService, "api/visit/v1"))

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return &statsResult{
			err: errors.Wrap(err, "could not send visit stats"),
		}
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logging.Println("can not read response body", err)
	}

	return &statsResult{
		nextDay: func(report *Report) {
			report.visitNextDay()
		},
	}
}

func CollectVisits(userId uuid.UUID, report *Report) stats {

	startDate := report.visitStartDate
	currentDay := startDate.Day()

	var visits []visit

	for {
		visits = append(visits, newVisit(userId, startDate))

		startDate = startDate.Add(2 * time.Hour)
		if startDate.Day() != currentDay {
			break
		}
	}
	return visitStats(visits)
}
