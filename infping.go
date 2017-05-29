// infping.go copyright Tor Hveem
// License: MIT

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Unable to read config file", err)
	}

	influxHost := viper.GetString("influx.host")
	influxPort := viper.GetString("influx.port")
	influxUser := viper.GetString("influx.user")
	influxPass := viper.GetString("influx.pass")
	influxDB := viper.GetString("influx.db")
	influxRetPolicy := viper.GetString("influx.policy")

	hosts := viper.GetStringSlice("hosts.hosts")

	u, err := url.Parse(fmt.Sprintf("http://%s:%s", influxHost, influxPort))
	if err != nil {
		log.Fatal("Unable to parse Influx Host/Port", err)
	}

	conf := client.HTTPConfig{
		Addr:     u.String(),
		Username: influxUser,
		Password: influxPass,
	}

	rawClient, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal("Failed to create Influx Client", err)
	}

	influxClient := NewInfluxClient(rawClient, influxDB, influxRetPolicy)

	dur, version, err := influxClient.Ping()
	if err != nil {
		log.Fatal("Unable to ping InfluxDB", err)
	}
	log.Printf("Pinged InfluxDB (version %s) in %v", version, dur)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Launching fping with hosts: %s", strings.Join(hosts, ", "))
	runAndRead(ctx, hosts, influxClient)
}
