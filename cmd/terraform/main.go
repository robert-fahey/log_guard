package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v2"
)

type LogCodes struct {
	AppName  string    `yaml:"appName"`
	LogCodes []LogCode `yaml:"logCodes"`
}

type LogCode struct {
	Code              int    `yaml:"code"`
	Level             string `yaml:"level"`
	Description       string `yaml:"description"`
	HumanReadableCode string `yaml:"humanReadableCode"`
}

func getLogCodesFromYaml() []LogCode {
	pwd, _ := os.Getwd()
	filename := filepath.Join(pwd, "log_codes.yaml")

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}

	var logCodes LogCodes
	err = yaml.Unmarshal(file, &logCodes)
	if err != nil {
		log.Fatalf("cannot unmarshal yaml: %v", err)
	}

	return logCodes.LogCodes
}

func getMetricNamesFromTf(tfFile string) []string {
	pwd, _ := os.Getwd()
	filename := filepath.Join(pwd, tfFile)

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("cannot read file: %v", err)
	}

	r := regexp.MustCompile(`jsonPayload.code=\\"(.+?)\\"`)
	matches := r.FindAllStringSubmatch(string(file), -1)

	var names []string
	for _, match := range matches {
		names = append(names, match[1])
	}

	return names
}

func checkCoverage(logCodes []LogCode, metricNames []string) []string {
	var uncovered []string

	for _, logCode := range logCodes {
		found := false
		for _, name := range metricNames {
			if name == logCode.HumanReadableCode {
				found = true
				break
			}
		}
		if !found {
			uncovered = append(uncovered, logCode.HumanReadableCode)
		}
	}

	return uncovered
}

func main() {

	logCodes := getLogCodesFromYaml()
	fmt.Println("Log codes from YAML:", logCodes)

	metricNames := getMetricNamesFromTf("main.tf")
	fmt.Println("Metric names from Terraform:")

	uncovered := checkCoverage(logCodes, metricNames)

	for _, logCode := range uncovered {
		fmt.Printf("Log code not covered in Terraform metrics: %v\n", logCode)
	}

	if len(uncovered) > 0 {
		os.Exit(1)
	}
}
