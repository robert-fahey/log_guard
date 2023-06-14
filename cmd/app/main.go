package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

type LogCode struct {
	Code              int    `yaml:"code"`
	Level             string `yaml:"level"`
	Description       string `yaml:"description"`
	HumanReadableCode string `yaml:"humanReadableCode"`
}

var appLogs []LogCode

func GenerateLogs() {
	appLogs = []LogCode{
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
}

func GetAppLogs() []LogCode {
	return appLogs
}

func main() {
	GenerateLogs()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	log.SetFormatter(&log.JSONFormatter{})

	for _, logCode := range GetAppLogs() {
		log.WithFields(log.Fields{
			"Code":              logCode.Code,
			"Level":             logCode.Level,
			"Description":       logCode.Description,
			"HumanReadableCode": logCode.HumanReadableCode,
		}).Info("Log entry")
	}
}
