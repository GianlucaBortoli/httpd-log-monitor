# loadgen
`loadgen` is a simple tool for simulating the load of incoming requests to an httpd web server.

It appends log lines to the `/tmp/access.log` log file with the following layout:
```
127.0.0.1 asd james [<date>] "GET <path> HTTP/1.0" <status_code> 123
```

The templated fields are filled as follows:
* `<date>` is always set to the time when the templated line is filled with data. This avoids appending
lines with a date in the past that would be ignored by `httpd-log-monitor`.<br>
The date layout is `02/Jan/2006:15:04:05 -0700`.
* `<path>` and `<status_code>` are randomly chosen at every iteration from predefined sets of
possible values.

## Run
The tool can be run even if `httpd-log-monitor` is not running. However, the main purpose of
`loadgen` is the manual testing of `httpd-log-monitor`. The common scenario is to open a terminal
to run `loadgen`:

```bash
$ cd loadgen/ # if you're not in this directory already
$ go run loadgen.go
```

and to open another terminal to run `httpd-log-monitor`. Please check [this section](../README.md#Run)
to know how to run it.

Please note that `httpd-log-monitor` won't start tailing if the log file (`/tmp/access.log` by default)
is not present on the machine. `loadgen` automatically creates the `/tmp/access.log` file if it's
not present, so please start `loadgen` before `httpd-log-monitor`.
