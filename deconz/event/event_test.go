package event

import (
	"errors"
	"os"
	"testing"
)

// examples from the xiaomi temp/hum/pressure sensor
const temperatureEventPayload = `{"e":"changed","id":"1","r":"sensors","state":{"lastupdated":"2018-03-08T19:35:24","temperature":2062},"t":"event"}`
const humidityEventPayload = `{"e":"changed","id":"2","r":"sensors","state":{"humidity":2985,"lastupdated":"2018-03-08T19:35:24"},"t":"event"}`
const pressureEventPayload = `{"e":"changed","id":"3","r":"sensors","state":{"lastupdated":"2018-03-08T19:35:24","pressure":993},"t":"event"}`

// xiaomi smoke detector
const smokeDetectorNoFireEventPayload = `{	"e": "changed",	"id": "5",	"r": "sensors",	"state": {	  "fire": false,	  "lastupdated": "2018-03-13T19:46:03",	  "lowbattery": false,	  "tampered": false	},	"t": "event"  }`

// xiaomi flood detector
const floodDetectorFloodDetectedEventPayload = `{ "e": "changed", "id": "6", "r": "sensors", "state": { "lastupdated": "2018-03-13T20:46:03", "lowbattery": false, "tampered": false, "water": true }, "t": "event"   }`

// xiaomi random switch "sensor"
const switchSensorEventPayload = `{	"e": "changed",	"id": "7",	"r": "sensors",	"state": {	  "buttonevent": 1000,	  "lastupdated": "2018-03-20T20:52:18"	},	"t": "event"  }  `

// some unknown json
const unknownEventPayload = `{"e": "hello my friend", "y": "Pænt go dag Hr 🦆"}`

type LookupImpl struct {
	Store map[int]string
}

func (l *LookupImpl) LookupType(i int) (string, error) {
	if t, ok := l.Store[i]; ok {
		return t, nil
	}
	return "", errors.New("not found")
}

func TestMain(m *testing.M) {
	TypeStore = &LookupImpl{Store: map[int]string{
		1: "ZHATemperature",
		2: "ZHAHumidity",
		3: "ZHAPressure",
		5: "ZHAFire",
		6: "ZHAWater",
		7: "ZHASwitch",
	}}
	os.Exit(m.Run())
}
func TestSmokeDetectorNoFireEvent(t *testing.T) {
	result, err := Parse([]byte(smokeDetectorNoFireEventPayload))
	if err != nil {
		t.Logf("unable to unmarshal smoke detector event: %s", err)
		t.FailNow()
	}

	smokeDetectorEvent, success := result.State.(ZHAFire)
	if !success {
		t.Log("unable to type assert smoke detector event")
		t.FailNow()
	}

	if smokeDetectorEvent.Fire != false {
		t.Fail()
	}
}

func TestUnknownEvent(t *testing.T) {
	_, err := Parse([]byte(unknownEventPayload))
	if err == nil {
		t.Fail()
	}
}

func TestFloodDetectorEvent(t *testing.T) {

	result, err := Parse([]byte(floodDetectorFloodDetectedEventPayload))
	if err != nil {
		t.Logf("Could not parse flood detector event: %s", err)
		t.FailNow()
	}

	floodEvent, success := result.State.(ZHAWater)
	if !success {
		t.Log("Unable to type assert floodevent")
		t.FailNow()
	}

	if !floodEvent.Water {
		t.Fail()
	}

}

func TestPressureEvent(t *testing.T) {

	result, err := Parse([]byte(pressureEventPayload))
	if err != nil {
		t.Logf("Could not parse pressure: %s", err)
		t.FailNow()
	}

	pressure, success := result.State.(ZHAPressure)
	if !success {
		t.Log("Coudl not assert to pressureevent")
		t.FailNow()
	}

	if pressure.Pressure != 993 {
		t.Fail()
	}
}

func TestTemperatureEvent(t *testing.T) {

	result, err := Parse([]byte(temperatureEventPayload))
	if err != nil {
		t.Logf("Could not parse temperature: %s", err)
		t.FailNow()
	}

	temp, success := result.State.(ZHATemperature)
	if !success {
		t.Logf("Could not assert to temperature event")
		t.FailNow()
	}

	if temp.Temperature != 2062 {
		t.Fail()
	}
}

func TestHumidityEvent(t *testing.T) {

	result, err := Parse([]byte(humidityEventPayload))
	if err != nil {
		t.Logf("Could not parse humidity: %s", err)
		t.FailNow()
	}

	humidity, success := result.State.(ZHAHumidity)
	if !success {
		t.Logf("unable assert humidity event")
		t.FailNow()
	}

	if humidity.Humidity != 2985 {
		t.Fail()
	}
}

func TestSwitchEvent(t *testing.T) {

	result, err := Parse([]byte(switchSensorEventPayload))
	if err != nil {
		t.Logf("Could not parse switch event: %s", err)
		t.FailNow()
	}

	s, success := result.State.(ZHASwitch)
	if !success {
		t.Logf("unable assert switch event")
		t.FailNow()
	}

	if s.Buttonevent != 1000 {
		t.Fail()
	}
}
