package deconz

import (
	"encoding/json"
	"fmt"
)

// Event represents a deconz sensor event
// {
//    "t": "event",
//    "e": "changed",
//    "r": "sensors",
//    "id": "7",
//    "state": { "presence": true }
// }
type Event struct {
	Type     string `json:"t"`
	Event    string `json:"e"`
	Resource string `json:"r"`
	ID       string `json:"id"`
}

// State is for embedding into event states
type State struct {
	Lastupdated string
}

// HumidityEvent represents a presure change
type HumidityEvent struct {
	Event
	State struct {
		State
		Humidity int
	}
}

// PressureEvent represents a presure change
type PressureEvent struct {
	Event
	State struct {
		State
		Pressure int
	}
}

// TemperatureEvent represents a temperature change
type TemperatureEvent struct {
	Event
	State struct {
		State
		Temperature int
	}
}

// Unmarshal decodes and returns apporiate event
func Unmarshal([]byte b) (interface{}, error) {
	dec.DisallowUnknownFields()

	var a TemperatureEvent
	err := dec.Decode(&a)
	if err == nil {
		return &a, nil
	}

	var b PressureEvent
	err = dec.Decode(&b)
	if err == nil {
		return &b, nil
	}

	var c HumidityEvent
	err = dec.Decode(&c)
	if err == nil {
		return &c, nil
	}

	return nil, fmt.Errorf("Unable to parse %s", err)
}
