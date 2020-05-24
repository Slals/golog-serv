package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logHandler)

	data := Log{
		Level:   "debug",
		Key:     "test_debug",
		Message: "This is a debug message",
	}
	dataByte, _ := json.Marshal(data)

	failedReq := httptest.NewRequest("PUT", "/logs", bytes.NewReader(dataByte))

	handler.ServeHTTP(rr, failedReq)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Bad http code, expected %d got %d", http.StatusNoContent, rr.Code)
	}
}

func TestBadLevelLogger(t *testing.T) {
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(logHandler)

	data := Log{
		Level:   "wrong_debug_level",
		Key:     "test_debug",
		Message: "This is a debug message",
	}
	dataByte, _ := json.Marshal(data)

	failedReq := httptest.NewRequest("PUT", "/logs", bytes.NewReader(dataByte))

	handler.ServeHTTP(rr, failedReq)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Bad http code, expected %d got %d", http.StatusBadRequest, rr.Code)
	}
}
