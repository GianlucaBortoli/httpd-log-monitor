package main

import (
	"flag"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/logmonitor"
)

func main() {
	file := flag.String("file", "/tmp/access.log", "The path to the log file")
	statsPeriod := flag.Duration("statsPeriod", 10*time.Second, "The period for displaying stats")
	statsK := flag.Int("statsK", 5, "The maximum number of stats to output every period")
	alertPeriod := flag.Duration("alertPeriod", 2*time.Minute, "The period for req/s alerting")
	alertThreshold := flag.Float64("alertThreshold", 10, "The req/s alert threshold")
	flag.Parse()

	m, err := logmonitor.New(*file, *alertPeriod, *statsPeriod, *statsK, *alertThreshold)
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}
	if err := m.Wait(); err != nil {
		panic(err)
	}
}
