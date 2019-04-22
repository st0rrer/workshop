package main

import (
	"github.com/st0rrer/workshop/web-service/src/config"
	reportHandler "github.com/st0rrer/workshop/web-service/src/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActivityHandler(t *testing.T) {

	config.GetConfig()

	badReq, err := http.NewRequest("POST", "/api/activity/v1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(reportHandler.GetReportHandler())

	handler.ServeHTTP(recorder, badReq)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status %v", status)
	}
}

func TestVisitHandler(t *testing.T) {

	config.GetConfig()

	req, err := http.NewRequest("POST", "/api/visit/v1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(reportHandler.GetReportHandler())

	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned unexpected status %v", status)
	}
}
