package report

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/st0rrer/workshop/client/src/config"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type activity struct {
	DataVer        int       `json:"dataVer"`
	UserId         uuid.UUID `json:"userId"`
	ActivityType   int       `json:"activityType"`
	StartTime      timestamp `json:"startTime"`
	EndTime        timestamp `json:"endTime"`
	StartLatitude  float64   `json:"startLatitude"`
	StartLongitude float64   `json:"startLongitude"`
	EndLatitude    float64   `json:"endLatitude"`
	EndLongitude   float64   `json:"endLongitude"`
}

func newActivity(userId uuid.UUID, startTime time.Time) activity {
	return activity{
		DataVer:        1,
		UserId:         userId,
		ActivityType:   RandomInt(1, 33, 2),
		StartTime:      timestamp(startTime),
		EndTime:        timestamp(startTime.Add(30 * time.Minute)),
		StartLatitude:  RandomLatitude(),
		StartLongitude: RandomLongitude(),
		EndLatitude:    RandomLatitude(),
		EndLongitude:   RandomLongitude(),
	}
}

type activityStats []activity

func (a activityStats) send() *statsResult {

	logger := log.New(os.Stdout, "activity#send ", 0)

	res, err := doPost(a, fmt.Sprintf("%v/%v", config.GetConfig().WebService, "api/activity/v1"))

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return &statsResult{
			err: errors.Wrap(err, "could not send activity stats"),
		}
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println("can not read response body", err)
	}

	return &statsResult{
		nextDay: func(report *Report) {
			report.activityNextDay()
		},
	}
}

func CollectActivities(userId uuid.UUID, report *Report) stats {

	startDate := report.activityStartDate
	currentDay := startDate.Day()

	var activities []activity

	for {
		activities = append(activities, newActivity(userId, startDate))

		startDate = startDate.Add(2 * time.Hour)
		if startDate.Day() != currentDay {
			break
		}
	}

	return activityStats(activities)
}
