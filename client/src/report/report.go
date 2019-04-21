package report

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
	"time"
)

var (
	collectors = []func(userId uuid.UUID, report *Report) stats{
		CollectVisits, CollectActivities,
	}
)

type Report struct {
	activityStartDate time.Time
	visitStartDate    time.Time
}

func NewReport() *Report {
	return &Report{
		activityStartDate: time.Date(2018, 1, 1, 1, 0, 0, 0, time.UTC),
		visitStartDate:    time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}

func (r *Report) activityNextDay() {
	r.activityStartDate = r.activityStartDate.AddDate(0, 0, 1)
}

func (r *Report) visitNextDay() {
	r.visitStartDate = r.visitStartDate.AddDate(0, 0, 1)
}

func (r *Report) ReportRecords(userId uuid.UUID) error {

	results := make(chan *statsResult, len(collectors))

	func(results chan *statsResult) {

		defer close(results)

		wg := sync.WaitGroup{}
		for _, collectStats := range collectors {

			collectStats := collectStats
			wg.Add(1)
			go func() {
				results <- collectStats(userId, r).send()
				wg.Done()
			}()
		}

		wg.Wait()

	}(results)

	var errors []string
	for result := range results {

		if result.err != nil {
			errors = append(errors, result.err.Error())
		} else {
			result.nextDay(r)
		}
	}

	if errors != nil {
		return fmt.Errorf(strings.Join(errors, "\n"))
	}
	return nil
}
