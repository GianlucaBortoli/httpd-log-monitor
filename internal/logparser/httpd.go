package logparser

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Songmu/axslogparser"
)

// HTTPd the apache httpd server log parser
type HTTPd struct {
	*axslogparser.Apache
}

// Line represents the parsed log line.
// See https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format for
// more information about the format.
type Line struct {
	RemoteHost    string
	RemoteLogName string
	User          string
	Date          time.Time
	Method        string
	Section       string
	Protocol      string
	StatusCode    int
	ContentLength int
}

// New returns an httpd log parser
func New() *HTTPd {
	return &HTTPd{&axslogparser.Apache{}}
}

// ParseLine takes a single log line and returns either its parsed version and an error
// in case the line is malformed or misses some required field (eg. the date)
func (p *HTTPd) ParseLine(line string) (*Line, error) {
	l, err := p.Parse(line)
	if err != nil {
		return nil, err
	}

	section, err := getSectionFromResource(l.RequestURI)
	if err != nil {
		return nil, err
	}
	return &Line{
		RemoteHost:    l.Host,
		RemoteLogName: l.RemoteLogname,
		User:          l.User,
		Date:          l.Time,
		Method:        l.Method,
		Section:       section,
		Protocol:      l.Protocol,
		StatusCode:    l.Status,
		ContentLength: int(l.Size),
	}, nil
}

// getSectionFromResource returns a section from a resource path.
// A section is defined as being what's before the second '/' in the resource section.
// Eg. the section for '/pages/create' is '/pages'.
// Returns an error if the path is the empty string
func getSectionFromResource(path string) (string, error) {
	parsed, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("cannot parse resource: %v", err)
	}
	if parsed.Path == "" {
		return "", fmt.Errorf("cannot get section on empty string path")
	}

	if !strings.HasPrefix(parsed.Path, "/") { // Reject paths that don't start with /
		return "", fmt.Errorf("cannot get section from path %s", parsed.Path)
	}
	// Remove leading "/" since I'm sure it's there
	stripped := strings.TrimLeft(parsed.Path, "/")
	// Split on middle "/"s and take the first
	split := strings.Split(stripped, "/")
	return "/" + split[0], nil
}
