package deconz

import (
	"fmt"

	"github.com/fasmide/deflux/deconz/event"
)

// SensorEvent is a sensor and a event embedded
type SensorEvent struct {
	*Sensor
	*event.Event
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

// WithSensor returns an sensor event with event and sensor embedded
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
