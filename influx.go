package main

import (
	"time"
	"github.com/influxdata/influxdb/client/v2"
)

// Client is a generic interface to write single-metric ping data
type Client interface {
	Write(point Point) error
	Ping() (time.Duration, string, error)
	Query(q client.Query) (*client.Response, error)
}

// NewInfluxClient creates a concrete InfluxDB Writer
func NewInfluxClient(client client.Client, db, retPolicy string) *InfluxClient {
	return &InfluxClient{
		influx:    client,
		db:        db,
		retPolicy: retPolicy,
	}
}

// InfluxClient implements the Client interface to provide a metrics client
// backed by InfluxDB
type InfluxClient struct {
	influx    client.Client
	db        string
	retPolicy string
}

// Ping calls Ping on the underlying influx client
func (i *InfluxClient) Ping() (time.Duration, string, error) {
	return i.influx.Ping(time.Second)
}

// Query calls Query on	the underlying influx client
func (i *InfluxClient) Query(q client.Query) (*client.Response, error) {
	return i.influx.Query(q)
}

// Write writes a single point to influx
func (i *InfluxClient) Write(point Point) error {
	var fields map[string]interface{}
	if point.Min != 0 && point.Avg != 0 && point.Max != 0 {
		fields = map[string]interface{}{
			"loss": point.LossPercent,
			"min":  point.Min,
			"avg":  point.Avg,
			"max":  point.Max,
		}
	} else {
		fields = map[string]interface{}{
			"loss": point.LossPercent,
		}
	}
	pt, err := client.NewPoint(
		"ping",
		map[string]string{
			"rx_host": point.RxHost,
			"tx_host": point.TxHost,
		},
		fields,
		time.Now())

	if err != nil {
		return err
	}

	batchConfig := client.BatchPointsConfig{Database: i.db, Precision: ""}
	bp, err := client.NewBatchPoints(batchConfig)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)
	if i.retPolicy != "" {
		bp.SetRetentionPolicy(i.retPolicy)
	}

	return i.influx.Write(bp)
}
