package deconz

import (
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

// some unknown json
const unknownEventPayload = `{"e": "hello my friend", "y": "Pænt go dag Hr 🦆"}`

func TestSmokeDetectorNoFireEvent(t *testing.T) {
	result, err := Unmarshal([]byte(smokeDetectorNoFireEventPayload))
	if err != nil {
		t.Logf("unable to unmarshal smoke detector event: %s", err)
		t.FailNow()
	}

	smokeDetectorEvent, success := result.(*SmokeDetectorEvent)
	if !success {
		t.Log("unable to type assert smoke detector event")
		t.FailNow()
	}

	if smokeDetectorEvent.State.Fire != false {
		t.Fail()
	}
}

func TestUnknownEvent(t *testing.T) {
	_, err := Unmarshal([]byte(unknownEventPayload))
	if err == nil {
		t.Fail()
	}
}

func TestFloodDetectorEvent(t *testing.T) {

	result, err := Unmarshal([]byte(floodDetectorFloodDetectedEventPayload))
	if err != nil {
		t.Logf("Could not parse flood detector event: %s", err)
		t.FailNow()
	}

	floodEvent, success := result.(*FloodEvent)
	if !success {
		t.Log("Unable to type assert floodevent")
		t.FailNow()
	}

	if !floodEvent.State.Water {
		t.Fail()
	}

}

func TestPressureEvent(t *testing.T) {

	result, err := Unmarshal([]byte(pressureEventPayload))
	if err != nil {
		t.Logf("Could not parse pressure: %s", err)
		t.FailNow()
	}

	pressure, success := result.(*PressureEvent)
	if !success {
		t.Log("Coudl not assert to pressureevent")
		t.FailNow()
	}

	if pressure.State.Pressure != 993 {
		t.Fail()
	}
}

func TestTemperatureEvent(t *testing.T) {

	result, err := Unmarshal([]byte(temperatureEventPayload))
	if err != nil {
		t.Logf("Could not parse temperature: %s", err)
		t.FailNow()
	}

	temp, success := result.(*TemperatureEvent)
	if !success {
		t.Logf("Could not assert to temperature event")
		t.FailNow()
	}

	if temp.State.Temperature != 2062 {
		t.Fail()
	}
}

func TestHumidityEvent(t *testing.T) {

	result, err := Unmarshal([]byte(humidityEventPayload))
	if err != nil {
		t.Logf("Could not parse humidity: %s", err)
		t.FailNow()
	}

	humidity, success := result.(*HumidityEvent)
	if !success {
		t.Logf("unable assert humidity event")
		t.FailNow()
	}

	if humidity.State.Humidity != 2985 {
		t.Logf("unexpected humidity value %d", humidity.State.Humidity)
		t.Fail()
	}
}
