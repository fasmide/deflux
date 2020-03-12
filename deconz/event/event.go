package event

import (
	"encoding/json"
	"fmt"
)

// TypeLookuper is the interface that we require to lookup types from id's
type TypeLookuper interface {
	LookupType(int) (string, error)
}

// Event represents a deconz sensor event
type Event struct {
	Type     string          `json:"t"`
	Event    string          `json:"e"`
	Resource string          `json:"r"`
	ID       int             `json:"id,string"`
	RawState json.RawMessage `json:"state"`
	State    interface{}
}

// Decoder is able to decode deCONZ events
type Decoder struct {
	TypeStore TypeLookuper
}

// Parse parses events from bytes
func (d *Decoder) Parse(b []byte) (*Event, error) {
	var e Event
	err := json.Unmarshal(b, &e)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal json: %s", err)
	}

	// If there is no state, dont try to parse it
	// TODO: figure out what to do with these
	//       some of them seems to be battery updates
	if e.Resource != "sensors" || len(e.RawState) == 0 {
		e.State = &EmptyState{}
		return &e, nil
	}

	err = e.ParseState(d.TypeStore)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal state: %s", err)
	}

	return &e, nil
}

// ParseState tries to unmarshal the appropriate state based
// on looking up the id though the TypeStore
func (e *Event) ParseState(tl TypeLookuper) error {

	t, err := tl.LookupType(e.ID)
	if err != nil {
		return fmt.Errorf("unable to lookup event id %d: %s", e.ID, err)
	}

	switch t {
	case "ZHAFire":
		var s ZHAFire
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHATemperature":
		var s ZHATemperature
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAPressure":
		var s ZHAPressure
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAHumidity":
		var s ZHAHumidity
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAWater":
		var s ZHAWater
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHASwitch":
		var s ZHASwitch
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAPresence":
		var s ZHAPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break	
	case "ClipPresence":
		var s ClipPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break	
	case "ZHALightLevel":
		var s ZHALightLevel
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAOpenClose":
		var s ZHAOpenClose
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "Daylight":
		var s Daylight
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
	default:
		err = fmt.Errorf("unable to unmarshal event state: %s is not a known type", t)
	}

	// err should continue to be null if everythings ok
	return err
}

// State is for embedding into event states
type State struct {
	Lastupdated string
}

// ZHAHumidity represents a presure change
type ZHAHumidity struct {
	State
	Humidity int
}

// Fields returns timeseries data for influxdb
func (z *ZHAHumidity) Fields() map[string]interface{} {
	return map[string]interface{}{
		"humidity": float64(z.Humidity) / 100,
	}
}

// ZHAPressure represents a presure change
type ZHAPressure struct {
	State
	Pressure int
}

// Fields returns timeseries data for influxdb
func (z *ZHAPressure) Fields() map[string]interface{} {
	return map[string]interface{}{
		"pressure": z.Pressure,
	}
}

// ZHATemperature represents a temperature change
type ZHATemperature struct {
	State
	Temperature int
}

// Fields returns timeseries data for influxdb
func (z *ZHATemperature) Fields() map[string]interface{} {
	return map[string]interface{}{
		"temperature": float64(z.Temperature) / 100,
	}
}

// ZHAPresence respresents a change from a presence sensor
type ClipPresence struct {
	State
	presence bool
}

// Fields returns timeseries data for influxdb
func (z *ClipPresence) Fields() map[string]interface{} {
	return map[string]interface{}{
		"presence": z.presence
	}
}

// ZHAPresence respresents a change from a presence sensor
type ZHAPresence struct {
	State
	presence bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAPresence) Fields() map[string]interface{} {
	return map[string]interface{}{
		"presence": z.presence
	}
}

// ZHAOpenClose respresents a change from a Door/Window open/close sensor
type ZHAOpenClose struct {
	State
	open bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAOpenClose) Fields() map[string]interface{} {
	return map[string]interface{}{
		"open": z.presence
	}
}

// ZHALightLevel respresents a change from a Lightlevel sensor
type ZHALightLevel struct {
	State
	dark bool,
	daylight bool,
	lightlevel int,
	lux int
}

// Fields returns timeseries data for influxdb
func (z *ZHALightLevel) Fields() map[string]interface{} {
	return map[string]interface{}{
		"dark": z.dark,
		"daylight": z.daylight,
		"lightlevel": z.lightlevel,
		"lux": z.lux
	}
}

// ZHAFire represents a change from a smoke detector
type ZHAFire struct {
	State
	Fire       bool
	Lowbattery bool
	Tampered   bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAFire) Fields() map[string]interface{} {
	return map[string]interface{}{
		"lowbattery": z.Lowbattery,
		"tampered":   z.Tampered,
		"fire":       z.Fire,
	}
}

// ZHASwitch represents a change from a button or switch
type ZHASwitch struct {
	State
	Buttonevent int
}

// Fields returns timeseries data for influxdb
func (z *ZHASwitch) Fields() map[string]interface{} {
	return map[string]interface{}{
		"buttonevent": z.Buttonevent,
	}
}

// ZHAWater respresents a change from a flood sensor
type ZHAWater struct {
	State
	Lowbattery bool
	Tampered   bool
	Water      bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAWater) Fields() map[string]interface{} {
	return map[string]interface{}{
		"lowbattery": z.Lowbattery,
		"tampered":   z.Tampered,
		"water":      z.Water,
	}
}

// ZHAPresence represents a presence seonsor
type ZHAPresence struct {
	State
	Daylight bool
	Status   int
}

// Fields returns timeseries data for influxdb
func (z *Daylight) Fields() map[string]interface{} {
	return map[string]interface{}{
		"daylight": z.Daylight,
		"status":   z.Status,
	}
}

// EmptyState is an empty struct used to indicate no state was parsed
type EmptyState struct{}
