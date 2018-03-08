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

// PresenceEvent represents a presence change
type PresenceEvent struct {
	Event
	State struct {
		Lastupdated string
		Presence    bool
	}
}

// TemperatureEvent represents a temperature change
type TemperatureEvent struct {
	Event
	State struct {
		Lastupdated string
		Temperature int
	}
}

// Parse takes a decoder and returns apporiate format
func Parse(dec *json.Decoder) (interface{}, error) {
	dec.DisallowUnknownFields()

	var a TemperatureEvent
	err := dec.Decode(&a)
	if err == nil {
		return &a, nil
	}

	var b PresenceEvent
	err = dec.Decode(&b)
	if err == nil {
		return &b, nil
	}

	return nil, fmt.Errorf("Unable to parse %s", err)
}
