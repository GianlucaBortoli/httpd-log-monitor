package logparser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	p := New()
	assert.NotNil(t, p)
	assert.IsType(t, &HTTPd{}, p)
}

func TestHTTPd_ParseLine(t *testing.T) {
	layout := "02/Jan/2006:15:04:05 -0700"
	date := "09/May/2018:16:00:39 +0000"
	dateTime, _ := time.Parse(layout, date)

	testCases := []struct {
		line      string
		expLine   *Line
		shouldErr bool
	}{
		{
			`127.0.0.1 asd james [` + date + `] "GET /report HTTP/1.0" 200 123`,
			&Line{
				RemoteHost:    "127.0.0.1",
				RemoteLogName: "asd",
				User:          "james",
				Date:          dateTime,
				Method:        "GET",
				Section:       "/report",
				Protocol:      "HTTP/1.0",
				StatusCode:    200,
				ContentLength: 123,
			},
			false,
		},
		{
			`127.0.0.1 - james [` + date + `] "GET /report HTTP/1.0" 200 123`,
			&Line{
				RemoteHost:    "127.0.0.1",
				RemoteLogName: "-",
				User:          "james",
				Date:          dateTime,
				Method:        "GET",
				Section:       "/report",
				Protocol:      "HTTP/1.0",
				StatusCode:    200,
				ContentLength: 123,
			},
			false,
		},
		{
			`127.0.0.1 - james [` + date + `] "GET /report/foo/bar HTTP/1.0" 200 123`,
			&Line{
				RemoteHost:    "127.0.0.1",
				RemoteLogName: "-",
				User:          "james",
				Date:          dateTime,
				Method:        "GET",
				Section:       "/report",
				Protocol:      "HTTP/1.0",
				StatusCode:    200,
				ContentLength: 123,
			},
			false,
		},
		{
			`127.0.0.1 - james [` + date + `] "GET http://example.com/report/foo/bar HTTP/1.0" 200 123`,
			&Line{
				RemoteHost:    "127.0.0.1",
				RemoteLogName: "-",
				User:          "james",
				Date:          dateTime,
				Method:        "GET",
				Section:       "/report",
				Protocol:      "HTTP/1.0",
				StatusCode:    200,
				ContentLength: 123,
			},
			false,
		},
		{
			"asd",
			nil,
			true,
		},
		{
			// Resource doesn't start with /
			`127.0.0.1 - james [` + date + `] "GET report HTTP/1.0" 200 123`,
			nil,
			true,
		},
		{
			// Missing date field
			`127.0.0.1 - james "GET /report HTTP/1.0" 200 123`,
			nil,
			true,
		},
		{
			// Date without timezone
			`127.0.0.1 - james [09/May/2018:16:00:39] "GET /report HTTP/1.0" 200 123`,
			nil,
			true,
		},
		{
			// Date without timezone and time
			`127.0.0.1 - james [09/May/2018] "GET /report HTTP/1.0" 200 123`,
			nil,
			true,
		},
		{
			// Some random byte slice converted to string
			string([]byte("0x1")),
			nil,
			true,
		},
	}

	p := New()
	for _, tt := range testCases {
		parsed, err := p.ParseLine(tt.line)
		assert.Equal(t, tt.shouldErr, err != nil)
		assert.EqualValues(t, tt.expLine, parsed)
	}
}

func TestGetSectionFromResource(t *testing.T) {
	testCases := []struct {
		resource   string
		expSection string
		shouldErr  bool
	}{
		{
			"/foo",
			"/foo",
			false,
		},
		{
			"/foo/bar",
			"/foo",
			false,
		},
		{
			"/foo/bar/baz",
			"/foo",
			false,
		},
		{
			"/",
			"/",
			false,
		},
		{
			"",
			"",
			false,
		},
		{
			"foo",
			"",
			true,
		},
		{
			"foo/bar",
			"",
			true,
		},
		{
			"foo/bar/baz",
			"",
			true,
		},
		{
			"85:asd//asd.asd",
			"",
			true,
		},
		{
			"http://example.com/foo",
			"/foo",
			false,
		},
		{
			"http://example.com/foo/bar",
			"/foo",
			false,
		},
		{
			"http://example.com/foo/bar/baz",
			"/foo",
			false,
		},
		{
			"http://example.com/",
			"/",
			false,
		},
		{
			"http://example.com",
			"",
			false,
		},
	}

	for _, tt := range testCases {
		section, err := getSectionFromResource(tt.resource)
		assert.Equal(t, tt.shouldErr, err != nil)
		assert.Equal(t, tt.expSection, section)
	}
}
