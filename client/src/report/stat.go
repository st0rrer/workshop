package report

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

var DefaultClient = &http.Client{
	Timeout: 10 * time.Second,
}

type stats interface {
	send() *statsResult
}

type statsResult struct {
	err     error
	nextDay func(r *Report)
}

func doPost(data interface{}, url string) (*http.Response, error) {

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	if err := json.NewEncoder(gz).Encode(data); err != nil {
		return nil, errors.Wrap(err, "could not serialize data")
	}

	if err := gz.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, errors.Wrap(err, "could not create request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	return DefaultClient.Do(req)
}
