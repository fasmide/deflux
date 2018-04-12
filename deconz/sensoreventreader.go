package deconz

import (
	"fmt"

	"github.com/fasmide/deflux/deconz/event"
)

// SensorLookup represents an interface for sensor lookup
type SensorLookup interface {
	LookupSensor(int) (*Sensor, error)
}

// EventReader interface
type EventReader interface {
	ReadEvent() (*event.Event, error)
}

// SensorEventReader reads events from an event.reader and returns SensorEvents
type SensorEventReader struct {
	lookup SensorLookup
	reader EventReader
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
