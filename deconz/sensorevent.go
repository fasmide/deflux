package deconz

import (
	"fmt"
	"strconv"

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
	f, ok := s.Event.State.(fielder)
	if !ok {
		return nil, nil, fmt.Errorf("this event (%T:%s) has no time series data", s.State, s.Name)
	}

	return map[string]string{"name": s.Name, "type": s.Sensor.Type, "id": strconv.Itoa(s.Event.ID)}, f.Fields(), nil
}
