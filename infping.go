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
	viper.SetDefault("influx.host", "localhost")
	viper.SetDefault("influx.port", "8086")
	viper.SetDefault("influx.user", "")
	viper.SetDefault("influx.pass", "")
	viper.SetDefault("influx.secure", false)
	viper.SetDefault("influx.db", "infping")

	viper.SetConfigName("infping")
	viper.AddConfigPath("/etc/")
	viper.AddConfigPath("/usr/local/etc/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Unable to read config file", err)
	}

	influxScheme := "https"
	if !viper.GetBool("influx.secure") {
		influxScheme = "http"
	}
	influxHost := viper.GetString("influx.host")
	influxPort := viper.GetString("influx.port")
	influxUser := viper.GetString("influx.user")
	influxPass := viper.GetString("influx.pass")
	influxDB := viper.GetString("influx.db")
	influxRetPolicy := viper.GetString("influx.policy")

	hosts := viper.GetStringSlice("hosts.hosts")

	u, err := url.Parse(fmt.Sprintf("%s://%s:%s", influsScheme, influxHost, influxPort))
	if err != nil {
		log.Fatal("Unable to build valid Influx URL", err)
	}

	conf := client.HTTPConfig{
		Addr:     u.String(),
		Username: influxUser,
		Password: influxPass,
	}

	rawClient, err := client.NewHTTPClient(conf)
	if err != nil {
		log.Fatal("Failed to create Influx client", err)
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
