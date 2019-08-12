package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	statusCodes = []int{
		200, 302, 404, 500,
	}

	pages = []string{
		"/foo/bar", "/some/page", "/a/b/c/d", "/x/z", "/",
	}
)

func randomStatusCode() int {
	c := rand.Int() % len(statusCodes)
	return statusCodes[c]
}

func randomPage() string {
	p := rand.Int() % len(pages)
	return pages[p]
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Example of valid log lime
	// `127.0.0.1 asd james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123`
	lineTpl := `127.0.0.1 asd james [%s] "GET %s HTTP/1.0" %d 123`

	for {
		now := time.Now().Format("02/Jan/2006:15:04:05 -0700")
		_, err = fmt.Fprintf(f, lineTpl+"\n", now, randomPage(), randomStatusCode())
		if err != nil {
			log.Fatal(f)
		}

		sl := rand.Int() % 100
		time.Sleep(time.Millisecond * time.Duration(50+sl)) // sleep for a minimum of 50ms up to 150ms
	}
}
