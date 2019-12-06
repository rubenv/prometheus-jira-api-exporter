package main

import (
	"log"

	jiraapiexporter "github.com/rubenv/prometheus-jira-api-exporter"
)

func main() {
	err := jiraapiexporter.Run()
	if err != nil {
		log.Fatal(err)
	}
}
