package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	randomCodes = []int{
		200, 404, 500, 302,
	}

	randomPages = []string{
		"/foo/bar", "/some/page", "/a/b/c/d", "/x/z", "/",
	}
)

func randomStatusCode() int {
	c := rand.Int() % len(randomCodes)
	return randomCodes[c]
}

func randomPage() string {
	p := rand.Int() % len(randomPages)
	return randomPages[p]
}

func main() {
	f, err := os.OpenFile("/tmp/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// `127.0.0.1 asd james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123`
	lineTpl := `127.0.0.1 asd james [%s] "GET %s HTTP/1.0" %d 123`

	for {
		sl := rand.Int() % 500
		now := time.Now().Format("02/Jan/2006:15:04:05 -0700")
		_, err = fmt.Fprintf(f, lineTpl+"\n", now, randomPage(), randomStatusCode())
		if err != nil {
			log.Fatal(f)
		}
		time.Sleep(time.Millisecond * time.Duration(50+sl))
	}
}
