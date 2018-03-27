package deconz

import (
	"errors"
	"fmt"

	"github.com/fasmide/deflux/deconz/event"
)

// SensorEvent is a sensor and a event embedded
type SensorEvent struct {
	*Sensor
	*event.Event
}

type fielder interface {
	Fields() map[string]interface{}
}

// Timeseries returns tags and fields for use in influxdb
func (s *SensorEvent) Timeseries() (map[string]string, map[string]interface{}, error) {
	fielder, ok := s.State.(fielder)
	if !ok {
		return nil, nil, errors.New("this event has no time series data")
	}

	return map[string]string{"name": s.Name, "type": s.Sensor.Type}, fielder.Fields(), nil
}

// SensorLookup represents an interface for sensor lookup
type SensorLookup interface {
	LookupSensor(int) (*Sensor, error)
}

// SensorEventReader reads events from an event.reader and returns SensorEvents
type SensorEventReader struct {
	lookup SensorLookup
	reader *event.Reader
}

// Read returns an sensor event with event and sensor embedded
func (s *SensorEventReader) Read() (*SensorEvent, error) {
	e, err := s.reader.ReadEvent()
	if err != nil {
		return nil, fmt.Errorf("unable to read event: %s", err)
	}

	sensor, err := s.lookup.LookupSensor(e.ID)
	if err != nil {
		return nil, fmt.Errorf("could not lookup sensor for id %d: %s", e.ID, err)
	}

	return &SensorEvent{Event: e, Sensor: sensor}, nil
}
