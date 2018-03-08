package deconz

import (
	"encoding/json"
	"strings"
	"testing"
)

// examples from the xiaomi temp/hum/pressure sensor
const temperatureEventPayload = `{"e":"changed","id":"1","r":"sensors","state":{"lastupdated":"2018-03-08T19:35:24","temperature":2062},"t":"event"}`
const humidityEventPayload = `{"e":"changed","id":"2","r":"sensors","state":{"humidity":2985,"lastupdated":"2018-03-08T19:35:24"},"t":"event"}`
const pressureEventPayload = `{"e":"changed","id":"3","r":"sensors","state":{"lastupdated":"2018-03-08T19:35:24","pressure":993},"t":"event"}`

func TestTemperatureEvent(t *testing.T) {

	dec := json.NewDecoder(strings.NewReader(temperatureEventPayload))

	result, err := Parse(dec)
	if err != nil {
		t.Logf("Could not parse temperature: %s", err)
		t.FailNow()
	}

	temp, success := result.(*TemperatureEvent)
	if !success {
		t.Logf("Could not assert to temperature event")
		t.FailNow()
	}
	t.Logf("We parsed %+v", temp)
	if temp.State.Temperature != 2062 {
		t.Fail()
	}
}
