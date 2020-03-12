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
	case "Daylight":
		var s Daylight
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAPresence":
		var s ZHAPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "CLIPPresence":
		var s CLIPPresence
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHALightLevel":
		var s ZHALightLevel
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAVibration":
		var s ZHAVibration
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHAOpenClose":
		var s ZHAOpenClose
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
	case "ZHACarbonMonoxide":
		var s ZHACarbonMonoxide
		err = json.Unmarshal(e.RawState, &s)
		e.State = &s
		break
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

// Daylight represents a change in daylight
type Daylight struct {
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

// ZHAPresence represents a presence Sensor
type ZHAPresence struct {
	State
	Presence bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAPresence) Fields() map[string]interface{} {
	return map[string]interface{}{
		"presence": z.Presence,
	}
}

// CLIPPresence represents a presence Sensor
type CLIPPresence struct {
	State
	Presence bool
}

// Fields returns timeseries data for influxdb
func (z *CLIPPresence) Fields() map[string]interface{} {
	return map[string]interface{}{
		"presence": z.Presence,
	}
}

// ZHALightLevel represents a LightLevel Sensor
type ZHALightLevel struct {
	State
	Dark bool
	Daylight bool
	LightLevel int32
	Lux int16
}

// Fields returns timeseries data for influxdb
func (z *ZHALightLevel) Fields() map[string]interface{} {
	return map[string]interface{}{
		"daylight": z.Daylight,
		"dark": z.Dark,
		"lightlevel": z.LightLevel,
		"lux": z.Lux,
	}
}

// ZHAVibration represents a Vibration Sensor
type ZHAVibration struct {
	State
	Vibration bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAVibration) Fields() map[string]interface{} {
	return map[string]interface{}{
		"vibration": z.Vibration,
	}
}

// ZHAVibration represents a Vibration Sensor
type ZHAOpenClose struct {
	State
	Open bool
}

// Fields returns timeseries data for influxdb
func (z *ZHAOpenClose) Fields() map[string]interface{} {
	return map[string]interface{}{
		"open": z.Open,
	}
}

// ZHACarbonMonoxide represents a CarbonMonoxide Sensor
type ZHACarbonMonoxide struct {
	State
	Carbonmonoxide bool
	Lowbattery bool
	Tampered bool
}

// Fields returns timeseries data for influxdb
func (z *ZHACarbonMonoxide) Fields() map[string]interface{} {
	return map[string]interface{}{
		"CO": z.Carbonmonoxide,
		"lowbattery": z.Lowbattery,
		"tampered": z.Tampered,
		}
}

// EmptyState is an empty struct used to indicate no state was parsed
type EmptyState struct{}
