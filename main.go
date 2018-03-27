package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"time"

	"github.com/fasmide/deflux/deconz"
	client "github.com/influxdata/influxdb/client/v2"
	yaml "gopkg.in/yaml.v2"
)

// YmlFileName is the filename
const YmlFileName = "deflux.yml"

// Configuration holds data for Deconz and influxdb configuration
type Configuration struct {
	Deconz           deconz.Config
	Influxdb         client.HTTPConfig
	InfluxdbDatabase string
}

func main() {
	config, err := loadConfiguration()
	if err != nil {
		log.Printf("no configuration could be found: %s", err)
		outputDefaultConfiguration()
		return
	}

	sensorChan, err := sensorEventChan(config.Deconz)
	if err != nil {
		panic(err)
	}

	// initial influx batch
	batch, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.InfluxdbDatabase,
		Precision: "s",
	})
	if err != nil {
		panic(err)
	}

	var timeout time.Timer

	influxdb, err := client.NewHTTPClient(config.Influxdb)
	if err != nil {
		panic(err)
	}

	for {

		select {
		case sensorEvent := <-sensorChan:
			tags, fields, err := sensorEvent.Timeseries()
			if err != nil {
				log.Printf("not adding event to influx batch: %s", err)
				continue
			}

			pt, err := client.NewPoint("deflux", tags, fields, time.Now())
			if err != nil {
				panic(err)
			}
			batch.AddPoint(pt)
			timeout.Reset(1 * time.Second)

		case <-timeout.C:
			// when timer fires: save batch points, initialize a new batch
			err := influxdb.Write(batch)
			if err != nil {
				panic(err)
			}

			log.Printf("Saved %d records to influxdb", len(batch.Points()))
			// influx batch point
			batch, err = client.NewBatchPoints(client.BatchPointsConfig{
				Database:  config.InfluxdbDatabase,
				Precision: "s",
			})
		}
	}
}

func sensorEventChan(c deconz.Config) (chan *deconz.SensorEvent, error) {
	// get an event reader from the API
	d := deconz.API{Config: c}
	reader, err := d.EventReader()
	if err != nil {
		return nil, err
	}

	// Dial the reader
	err = reader.Dial()
	if err != nil {
		return nil, err
	}

	// create a new reader, embedding the event reader
	sensorEventReader := d.SensorEventReader(reader)
	channel := make(chan *deconz.SensorEvent)

	go func() {
		for {
			e, err := sensorEventReader.Read()
			if err != nil {
				log.Printf("Error received from sensoreventreader: %s", err)
				close(channel)
				return
			}
			channel <- e
		}
	}()

	return channel, nil
}

func loadConfiguration() (*Configuration, error) {
	data, err := readConfiguration()
	if err != nil {
		return nil, fmt.Errorf("could not read configuration: %s", err)
	}

	var config Configuration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse configuration: %s", err)
	}
	return &config, nil
}

// readConfiguration tries to read pwd/deflux.yml or /etc/deflux.yml
func readConfiguration() ([]byte, error) {
	// first try to load ${pwd}/deflux.yml
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("unable to get current work directory: %s", err)
	}

	pwdPath := path.Join(pwd, YmlFileName)
	data, pwdErr := ioutil.ReadFile(pwdPath)
	if pwdErr == nil {
		log.Printf("Using configuration %s", pwdPath)
		return data, nil
	}

	// if we reached this code, we where unable to read a "local" Configuration
	// try from /etc/deflux.yml
	etcPath := path.Join("/etc", YmlFileName)
	data, etcErr := ioutil.ReadFile(etcPath)
	if etcErr != nil {
		return nil, fmt.Errorf("\n%s\n%s", pwdErr, etcErr)
	}

	log.Printf("Using configuration %s", etcPath)
	return data, nil
}

func outputDefaultConfiguration() {

	c := defaultConfiguration()
	yml, err := yaml.Marshal(c)
	if err != nil {
		log.Fatalf("unable to generate default configuration: %s", err)
	}

	u, err := url.Parse(c.Deconz.Addr)
	if err == nil {
		// try to pair with deCONZ
		apikey, err := deconz.Pair(*u)
		if err != nil {
			log.Printf("unable to pair with deconz: %s, please fill out APIKey manually", err)
		}
		c.Deconz.APIKey = string(apikey)
	}

	log.Printf("Outputting default configuration, save this to /etc/deflux.yml")
	// to stdout
	fmt.Print(string(yml))
}

func defaultConfiguration() *Configuration {
	// this is the default configuration
	c := Configuration{
		Deconz: deconz.Config{
			Addr:   "http://127.0.0.1:8080/",
			APIKey: "change me",
		},
		Influxdb: client.HTTPConfig{
			Addr:      "http://127.0.0.1:8086/",
			Username:  "change me",
			Password:  "change me",
			UserAgent: "Deflux",
		},
		InfluxdbDatabase: "deconz",
	}

	// lets see if we are able to discover a gateway, and overwrite parts of the
	// default congfiguration
	discovered, err := deconz.Discover()
	if err != nil {
		log.Printf("discovery of deconz gateway failed: %s, please fill configuration manually..", err)
		return &c
	}

	// TODO: discover is acturlly a slice of multiple discovered gateways,
	// but for now we use only the first available
	deconz := discovered[0]
	addr := url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", deconz.InternalIPAddress, deconz.InternalPort),
		Path:   "/api",
	}
	c.Deconz.Addr = addr.String()

	return &c
}
