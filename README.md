# infping
Parse fping output, store results in influxdb

See blog post for more info https://hveem.no/visualizing-latency-variance-with-grafana

Simply run :
```go
go get github.com/influxdata/influxdb/client
go get -d github.com/imkwx/infping
```
Edit the config.toml file and add the retention policy in infping.go then

```go
go install github.com/imkwx/infping
go run bin/infping
```
