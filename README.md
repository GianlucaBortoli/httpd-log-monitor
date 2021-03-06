# httpd-log-monitor
An httpd log monitor with console alerting.

## Description
`httpd-log-monitor` is a monitor for the Apache httpd web server. It scrapes its log file (see
[here](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format) for the log format)
to compute some metrics, such as the requests second handled and the top visited sections of the web
site.

The monitor also prints an alert on the console if the requests/second threshold is reached as
well as another message when the alert is resolved.

## Build
The project uses go modules and requires go `>= 1.12`. To build the binary run
```bash
$ make build
```

The executable will be placed in `./bin` and has the race detector enabled.

## Unit tests
One-line assertions in unit-tests are provided by https://github.com/stretchr/testify. It's a handy
wrapper of the golang standard `testing` package.

Unit tests can be run (with the race detector enabled) typing 
```bash
$ make test
```

The overall code coverage is also outputted.

## Manual tests
Please refer to [this page](./loadgen/).

## Run
The `./bin/httpd-log-monitor` binary can be run without any further setting as follows:
```bash
$ ./bin/httpd-log-monitor
```

However, it's possible to override the default settings via command line parameters.
All the possible parameters can be found running the binary with the `-h` flag:
```bash
$ ./bin/httpd-log-monitor -h
Usage of ./bin/httpd-log-monitor:
  -alertPeriod duration
    	The length of the period for computing the request rate metric used for alerting about high traffic conditions (default 2m0s)
  -alertThreshold float
    	The threshold on the request rate metric for alerting about high traffic conditions (default 10)
  -logFile string
    	The path to the log file (default "/tmp/access.log")
  -statsK int
    	The maximum number of values to output when displaying topK metrics (eg. sections) (default 5)
  -statsPeriod duration
    	The length of the period for computing all the metrics and displaying them on the console (default 10s)
```

The process exits immediately with an error message if the log file doesn't exist.

## Metrics
The following metrics are collected from the log file:
* Rate of requests: the number requests per second. It indicates the load the web server is facing.
* Rate of errors: the number of status codes indicating an error. It indicates the health status of
the content being served or the presence of misbehaving clients.<br>
For example, a lot of `404`s can signal a wrong link in a web page or a web crawler that is scraping
all possible paths for the web server.
* TopK sections: the top `K` visited sections.
* TopK status codes: the top `K` status codes returned.
* TopK users: the top `K` users who did the request.

## Design decisions
Some design decisions and trade-offs have been made during the development of this tool.
Here is a list of the main ones divided by topics.

* Log line parsing:
    * Accepted line format is defined [here](https://www.w3.org/Daemon/User/Config/Logging.html#common-logfile-format).
    Eg:<br>
    `127.0.0.1 - james [09/May/2018:16:00:39 +0000] "GET /report HTTP/1.0" 200 123`
    * Malformed lines are gracefully handled but will be completely ignored.
* Log file tail:
    * The date in the log line is used to skip old lines. This is important when the tool is run against
    a file that already has some content (eg. when the web server is already running). The tool starts
    collecting data from the moment it’s run not considering old and stale data.
    * The behaviour should be the same of `tail -f`.
    * Survive file truncation during during the tailing process. In real-world examples, it is very
    common for log files to be truncated at some point. This is what log rotation tools usually do,
    hence this is handled by `http-log-monitor`.
* Collected metrics:
    * Sections of the web site with the most hits (topK, with configurable `K` via CLI parameter).
    * Average requests per second.
    * The metrics are handled as batches of a certain time length (configurable via CLI parameter).
* Alerts:
    * When the average of req/sec in the alerting time frame goes reaches the threshold (configurable
    via CLI parameter) a "high traffic" alert message is printed to che console.
    * If an alert fired, another message is printed to the console when the value goes below the
    threshold. This means that the alert is now resolved.


## Improvements
Given the decisions listed in the previous section, it's possible to think of possible improvements
that can enhance both performance, stability and maintainability.

* Save somewhere the last known position in the log file so the tool can start tailing from that point
onwards instead of always starting from scratch. The timestamp check may still be needed, but we could
avoid parsing many log lines just to skip them.
* Handle a gentle shutdown when either SIGINT or SIGTERM is sent to the `httpd-log-monitor` process.
This will ensure the tool cleans up after itself when exiting. For example, the tailer should remove
the inotify watches added by the tail package, since the Linux kernel may not automatically remove
inotify watches after the process exits (see [here](https://godoc.org/github.com/hpcloud/tail#Tail.Cleanup)
for more information).
* Introduce a common interface that every metric should implement. This will standardize the lifecycle
of every metric and it will make the code more maintainable if the number of metrics grows. Moreover,
it would allows to implement and create alerts on arbitrary metrics and not just on some specific ones.
* If the number of metrics for the `Manager` grows, it would need some change to keep track of them
in a more handy way. The current implementation has an event loop that serializes the access to all
the metrics objects since their state is kept in memory. This could become a bottleneck in case of
a huge increase in the number of metrics and/or data points to ingest. One possible solution is to
store the state of all metrics into another component (eg. an external database) with enough consistency
and performance guarantees. Given the limited scope of this project, I don't think this can become an
issue since the load on the `Manager` isn't too high at the moment.

## Known limitations
The library used for tailing the log file (https://github.com/hpcloud/tail) has known problems under
Windows (see [here](https://github.com/hpcloud/tail/labels/Windows)).
This tool has been tested on Ubuntu Linux, but Apple's macOS should be fine as well.
