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

// WithSensor returns an sensor event with event and sensor embedded
func WithSensor(e *event.Event, s SensorLookup) (*SensorEvent, error) {
	sensor, err := s.LookupSensor(e.ID)
	if err != nil {
		return nil, fmt.Errorf("could not lookup sensor for id %d: %s", e.ID, err)
	}

	return &SensorEvent{Event: e, Sensor: sensor}, nil
}
