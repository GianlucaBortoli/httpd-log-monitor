package main

import (
	"flag"

	"github.com/cog-qlik/httpd-log-monitor/pkg/monitor"
)

func main() {
	file := flag.String("file", "/tmp/access.log", "The path to the log file")
	flag.Parse()

	m := monitor.New(*file)

	if err := m.Start(); err != nil {
		panic(err)
	}
	if err := m.Wait(); err != nil {
		panic(err)
	}
}
