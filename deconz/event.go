package deconz

import (
	"bytes"
	"encoding/json"
	"errors"
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

// FloodEvent respresents an event from a flood sensor
type FloodEvent struct {
	Event
	State struct {
		State
		Lowbattery bool
		Tampered   bool
		Water      bool
	}
}

// SmokeDetectorEvent resporesents an event from a smoke detector
type SmokeDetectorEvent struct {
	Event
	State struct {
		State
		Fire       bool
		Lowbattery bool
		Tampered   bool
	}
}

// Unmarshal decodes and returns apporiate event
// TODO: This should at least be made more robust, for example
// passing {e: "something"} will make unmarshal return an somewhat
// faulty TemperatureEvent...
// this is not really what we want when dealing with
// something like a smoke detector
func Unmarshal(payload []byte) (interface{}, error) {

	// in order to use DisallowUnknownFields we cannot use
	// the typical json.Unmarshal we must create a decoder
	// in order to pass the same payload multiple times, we
	// create this buffer and write the payload again as we
	// along...
	buf := bytes.NewBuffer(payload)

	dec := json.NewDecoder(buf)
	dec.DisallowUnknownFields()

	var a TemperatureEvent
	err := dec.Decode(&a)
	if err == nil {
		return &a, nil
	}

	// if the above failed, we must rewrite the same byte slice
	// otherwise the decoder will just read EOF
	buf.Write(payload)

	var b PressureEvent
	err = dec.Decode(&b)
	if err == nil {
		return &b, nil
	}

	buf.Write(payload)

	var c HumidityEvent
	err = dec.Decode(&c)
	if err == nil {
		return &c, nil
	}

	buf.Write(payload)

	var d FloodEvent
	err = dec.Decode(&d)
	if err == nil {
		return &d, nil
	}

	buf.Write(payload)

	var e SmokeDetectorEvent
	err = dec.Decode(&e)
	if err == nil {
		return &e, nil
	}

	return nil, errors.New("payload did not match any of our event types")
}
