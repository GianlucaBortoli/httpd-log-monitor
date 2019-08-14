package main

import (
	"flag"
	"log"
	"time"

	"github.com/cog-qlik/httpd-log-monitor/pkg/logmonitor"
)

var (
	logFile        = flag.String("logFile", "/tmp/access.log", "The path to the log file")
	statsPeriod    = flag.Duration("statsPeriod", 10*time.Second, "The length of the period for computing all the metrics and displaying them on the console")
	statsK         = flag.Int("statsK", 5, "The maximum number of values to output when displaying topK metrics (eg. sections)")
	alertPeriod    = flag.Duration("alertPeriod", 2*time.Minute, "The length of the period for computing the request rate metric used for alerting about high traffic conditions")
	alertThreshold = flag.Float64("alertThreshold", 10, "The threshold on the request rate metric for alerting about high traffic conditions")
)

func main() {
	flag.Parse()

	m, err := logmonitor.New(*logFile, *alertPeriod, *statsPeriod, *statsK, *alertThreshold)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Start(); err != nil {
		log.Fatal(err)
	}

	if err = m.Wait(); err != nil {
		log.Fatal(err)
	}
}
