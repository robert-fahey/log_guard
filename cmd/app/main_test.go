package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

type LogCodes struct {
	LogCodes []LogCode `yaml:"logCodes"`
}

func ReadLogCodesFile(filename string) (*LogCodes, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	codes := &LogCodes{}
	err = yaml.Unmarshal(buf, codes)
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func TestLogCodesCoverage(t *testing.T) {
	pwd, _ := os.Getwd()
	filename := filepath.Join(pwd, "../../log_codes.yaml")

	logCodes, err := ReadLogCodesFile(filename)
	if err != nil {
		log.Fatalf("failed to read log codes file: %v", err)
	}
	GenerateLogs()
	appLogs := GetAppLogs()

	for _, code := range logCodes.LogCodes {
		found := false
		for _, appLog := range appLogs {
			if code.Code == appLog.Code &&
				code.Level == appLog.Level &&
				code.Description == appLog.Description &&
				code.HumanReadableCode == appLog.HumanReadableCode {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("log code not covered in application logs: %+v", code)
		}
	}
}
