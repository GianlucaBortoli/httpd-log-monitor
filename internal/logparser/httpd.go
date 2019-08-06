package logparser

import (
	"time"

	"github.com/Songmu/axslogparser"
)

type HTTPd struct {
	parser *axslogparser.Apache
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
	Path          string
	Protocol      string
	StatusCode    int
	ContentLength int
}

func New() *HTTPd {
	return &HTTPd{&axslogparser.Apache{}}
}

func (p *HTTPd) ParseLine(line string) (*Line, error) {
	l, err := p.parser.Parse(line)
	if err != nil {
		return nil, err
	}

	return &Line{
		RemoteHost:    l.Host,
		RemoteLogName: l.RemoteLogname,
		User:          l.User,
		Date:          l.Time,
		Method:        l.Method,
		Path:          l.RequestURI,
		Protocol:      l.Protocol,
		StatusCode:    l.Status,
		ContentLength: int(l.Size),
	}, nil
}
