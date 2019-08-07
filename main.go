package main

import (
	"flag"

	"github.com/cog-qlik/httpd-log-monitor/pkg/logmonitor"
)

func main() {
	file := flag.String("file", "/tmp/access.log", "The path to the log file")
	flag.Parse()

	m := logmonitor.New(*file)

	if err := m.Start(); err != nil {
		panic(err)
	}
	if err := m.Wait(); err != nil {
		panic(err)
	}
}
