package main

import (
	"flag"
	"fmt"

	"github.com/cog-qlik/httpd-log-monitor/internal/tailer"
)

func main() {
	file := flag.String("file", "/tmp/access.log", "The log file path to tail")
	flag.Parse()

	t := tailer.New(*file)
	lines, err := t.Start()
	if err != nil {
		panic(err)
	}

	for l := range lines {
		fmt.Println(l.Text)
	}
	if err := t.Wait(); err != nil {
		panic(err)
	}
}
