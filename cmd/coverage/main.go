package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// LogCode struct represents the structure of log codes in the YAML file
type LogCode struct {
	Code              int    `yaml:"code"`
	Level             string `yaml:"level"`
	Description       string `yaml:"description"`
	HumanReadableCode string `yaml:"humanReadableCode"`
}

// LogEntry struct represents the structure of log entries from the standard input
type LogEntry struct {
	Code              int    `json:"Code"`
	Description       string `json:"Description"`
	HumanReadableCode string `json:"HumanReadableCode"`
	Level             string `json:"Level"`
	Msg               string `json:"msg"`
	LogLevel          string `json:"level"`
	Time              string `json:"time"`
}

// LogCodes struct is a container for a slice of LogCode structs
type LogCodes struct {
	LogCodes []LogCode `yaml:"logCodes"`
}

// getLogCodesFromYaml function reads log codes from a YAML file
func getLogCodesFromYaml() ([]LogCode, error) {
	filename := os.Getenv("LOG_CODES_YAML")
	if filename == "" {
		return nil, fmt.Errorf("environment variable LOG_CODES_YAML not set")
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %v", err)
	}

	var logCodes LogCodes
	err = yaml.Unmarshal(file, &logCodes)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal yaml: %v", err)
	}

	return logCodes.LogCodes, nil
}

// getLogsFromStdin function reads log entries from the standard input
func getLogsFromStdin() ([]LogEntry, error) {
	reader := bufio.NewReader(os.Stdin)

	var logs []LogEntry
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error occurred while reading line: %v", err)
		}

		var logEntry LogEntry
		if err := json.Unmarshal(line, &logEntry); err != nil {
			return nil, fmt.Errorf("error occurred while parsing line: %v. Error: %v", string(line), err)
		}
		logs = append(logs, logEntry)
	}

	return logs, nil
}

// Helper function to check if a log code is covered
func isLogCodeCovered(logCode LogCode, logs []LogEntry) bool {
	for _, log := range logs {
		if log.HumanReadableCode == logCode.HumanReadableCode {
			return true
		}
	}
	return false
}

// checkCoverage function checks which log codes are not covered by the logs
func checkCoverage(logCodes []LogCode, logs []LogEntry) []LogCode {
	var uncovered []LogCode
	for _, logCode := range logCodes {
		if !isLogCodeCovered(logCode, logs) {
			uncovered = append(uncovered, logCode)
		}
	}
	return uncovered
}

func main() {
	log.SetOutput(os.Stderr)

	// Get log codes from YAML file
	logCodes, err := getLogCodesFromYaml()
	if err != nil {
		log.Fatalf("Failed to get log codes from YAML: %v", err)
	}

	// Get log entries from standard input
	logs, err := getLogsFromStdin()
	if err != nil {
		log.Fatalf("Failed to get logs from standard input: %v", err)
	}

	// Check which log codes are not covered by the logs
	uncovered := checkCoverage(logCodes, logs)

	log.SetOutput(os.Stdout)

	totalCodes := len(logCodes)
	uncoveredCodes := len(uncovered)
	coveredCodes := totalCodes - uncoveredCodes
	coveragePercent := float64(coveredCodes) / float64(totalCodes) * 100

	// ASCII Art
	log.Println(`
█    ████▄   ▄▀    ▄▀    ▄   ██   █▄▄▄▄ ██▄   
█    █   █ ▄▀    ▄▀       █  █ █  █  ▄▀ █  █  
█    █   █ █ ▀▄  █ ▀▄  █   █ █▄▄█ █▀▀▌  █   █ 
███▄ ▀████ █   █ █   █ █   █ █  █ █  █  █  █  
    ▀       ███   ███  █▄ ▄█    █   █   ███▀  
                        ▀▀▀    █   ▀          
                              ▀                        
	`)

	if uncoveredCodes == 0 {
		log.Println("All log codes are covered in application logs. Coverage: 100%")
	} else {
		for _, logCode := range uncovered {
			log.Printf("Log code not covered in application logs: %v\n", logCode.HumanReadableCode)
		}
		log.Printf("Coverage: %.2f%%", coveragePercent)
	}
}
