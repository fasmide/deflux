package deconz

import (
	"strconv"
	"testing"

	"github.com/fasmide/deflux/deconz/event"
)

const smokeDetectorNoFireEventPayload = `{	"e": "changed",	"id": "5",	"r": "sensors",	"state": {	  "fire": false,	  "lastupdated": "2018-03-13T19:46:03",	  "lowbattery": false,	  "tampered": false	},	"t": "event"  }`

type testLookup struct {
}

func (t *testLookup) LookupSensor(i int) (*Sensor, error) {
	return &Sensor{Name: "Test Sensor", Type: "ZHAFire"}, nil
}

func (t *testLookup) LookupType(i int) (string, error) {
	return "ZHAFire", nil
}

type testReader struct {
}

func (t *testReader) ReadEvent() (*event.Event, error) {
	d := event.Decoder{TypeStore: &testLookup{}}
	return d.Parse([]byte(smokeDetectorNoFireEventPayload))
}
func TestSensorEventReader(t *testing.T) {

	r := SensorEventReader{lookup: &testLookup{}, reader: &testReader{}}

	e, err := r.Read()
	if err != nil {
		t.Fail()
	}
	if strconv.Itoa(e.Event.ID) != "5" {
		t.Fail()
	}
	tags, fields, err := e.Timeseries()
	if err != nil {
		t.Logf(err.Error())
		t.FailNow()
	}
	if tags["name"] != "Test Sensor" {
		t.Fail()
	}
	if tags["id"] != "5" {
		t.Fail()
	}

	if fields["fire"] != false {
		t.Fail()
	}

}
