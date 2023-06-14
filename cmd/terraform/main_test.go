package main

import (
	"reflect"
	"testing"
)

func TestCheckCoverage(t *testing.T) {
	logCodes := []LogCode{
		{Code: 100, Level: "ERROR", Description: "Unexpected null pointer encountered.", HumanReadableCode: "ERROR_NULL_POINTER"},
		{Code: 200, Level: "WARN", Description: "Data format mismatch, falling back to default.", HumanReadableCode: "WARN_DATA_FORMAT_MISMATCH"},
	}
	metricNames := []string{"ERROR_NULL_POINTER"}

	expected := []string{"WARN_DATA_FORMAT_MISMATCH"}
	actual := checkCoverage(logCodes, metricNames)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
