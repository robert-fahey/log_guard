package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func getLogsFromReader(r io.Reader) ([]LogEntry, error) {
	var logs []LogEntry
	dec := json.NewDecoder(r)
	err := dec.Decode(&logs)
	if err != nil {
		return nil, fmt.Errorf("error occurred while parsing line: %v", err)
	}

	return logs, nil
}

func TestGetLogsFromReader(t *testing.T) {
	jsonLogs := `[
		{
			"Code": 100,
			"Description": "Unexpected null pointer encountered.",
			"HumanReadableCode": "ERROR_NULL_POINTER",
			"Level": "ERROR",
			"LogLevel": "info",
			"Msg": "Log entry",
			"Time": "2023-06-14T08:26:13+02:00"
		},
		{
			"Code": 200,
			"Description": "Data format mismatch, falling back to default.",
			"HumanReadableCode": "WARN_DATA_FORMAT_MISMATCH",
			"Level": "WARN",
			"LogLevel": "info",
			"Msg": "Log entry",
			"Time": "2023-06-14T08:26:13+02:00"
		}
	]
	`
	reader := bytes.NewBufferString(jsonLogs)

	logs, err := getLogsFromReader(reader)
	if err != nil {
		t.Fatalf("Failed to get logs from reader: %v", err)
	}

	expected := []LogEntry{
		{
			Code:              100,
			Description:       "Unexpected null pointer encountered.",
			HumanReadableCode: "ERROR_NULL_POINTER",
			Level:             "ERROR",
			Msg:               "Log entry",
			Time:              "2023-06-14T08:26:13+02:00",
		},
		{
			Code:              200,
			Description:       "Data format mismatch, falling back to default.",
			HumanReadableCode: "WARN_DATA_FORMAT_MISMATCH",
			Level:             "WARN",
			Msg:               "Log entry",
			Time:              "2023-06-14T08:26:13+02:00",
		},
		// add more expected log entries as required
	}

	if !reflect.DeepEqual(logs, expected) {
		t.Errorf("getLogsFromReader() = %v, want %v", logs, expected)
	}
}

func TestGetLogCodesFromYaml(t *testing.T) {
	yamlLogs := `
	- Code: 100
	  Description: Unexpected null pointer encountered.
	  HumanReadableCode: ERROR_NULL_POINTER
	  Level: ERROR
	- Code: 200
	  Description: Data format mismatch, falling back to default.
	  HumanReadableCode: WARN_DATA_FORMAT_MISMATCH
	  Level: WARN
	`
	appFS := afero.NewMemMapFs() // use afero's in-memory filesystem for testing
	err := afero.WriteFile(appFS, "log_codes.yaml", []byte(yamlLogs), 0644)
	if err != nil {
		t.Fatalf("Failed to write to in-memory fs: %v", err)
	}

	logCodes, err := getLogCodesFromYaml()
	if err != nil {
		t.Fatalf("Failed to get logs from reader: %v", err)
	}

	expected := []LogCode{
		{
			Code:              100,
			Level:             "ERROR",
			Description:       "Unexpected null pointer encountered.",
			HumanReadableCode: "ERROR_NULL_POINTER",
		},
		{
			Code:              200,
			Level:             "WARN",
			Description:       "Data format mismatch, falling back to default.",
			HumanReadableCode: "WARN_DATA_FORMAT_MISMATCH",
		},
		{
			Code:              300,
			Level:             "INFO",
			Description:       "Database connection successfully established.",
			HumanReadableCode: "INFO_DB_CONNECTION_ESTABLISHED",
		},
		{
			Code:              400,
			Level:             "DEBUG",
			Description:       "Starting API endpoint health check.",
			HumanReadableCode: "DEBUG_HEALTH_CHECK",
		},
		{
			Code:              500,
			Level:             "TRACE",
			Description:       "Beginning detailed transaction trace.",
			HumanReadableCode: "TRACE_TRANSACTION",
		},
	}

	if !reflect.DeepEqual(logCodes, expected) {
		t.Errorf("getLogCodesFromYaml() = %v, want %v", logCodes, expected)
	}
}

func TestCheckCoverage(t *testing.T) {
	logCodes := []LogCode{
		{
			Code:              100,
			Level:             "ERROR",
			Description:       "Unexpected null pointer encountered.",
			HumanReadableCode: "ERROR_NULL_POINTER",
		},
		{
			Code:              200,
			Level:             "WARN",
			Description:       "Data format mismatch, falling back to default.",
			HumanReadableCode: "WARN_DATA_FORMAT_MISMATCH",
		},
		{
			Code:              300,
			Level:             "INFO",
			Description:       "Database connection successfully established.",
			HumanReadableCode: "INFO_DB_CONNECTION_ESTABLISHED",
		},
	}

	logs := []LogEntry{
		{
			Code:              100,
			Description:       "Unexpected null pointer encountered.",
			HumanReadableCode: "ERROR_NULL_POINTER",
			Level:             "ERROR",
			Msg:               "Log entry",
			LogLevel:          "info",
			Time:              "2023-06-14T08:26:13+02:00",
		},
		{
			Code:              200,
			Description:       "Data format mismatch, falling back to default.",
			HumanReadableCode: "WARN_DATA_FORMAT_MISMATCH",
			Level:             "WARN",
			Msg:               "Log entry",
			LogLevel:          "info",
			Time:              "2023-06-14T08:26:13+02:00",
		},
	}

	uncovered := checkCoverage(logCodes, logs)

	expected := []LogCode{
		{
			Code:              300,
			Level:             "INFO",
			Description:       "Database connection successfully established.",
			HumanReadableCode: "INFO_DB_CONNECTION_ESTABLISHED",
		},
	}

	if !reflect.DeepEqual(uncovered, expected) {
		t.Errorf("checkCoverage() = %v, want %v", uncovered, expected)
	}
}
