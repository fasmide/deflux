package deconz

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/runger1101001/deflux/deconz/event"
)

// API represents the deCONZ rest api
type API struct {
	Config      Config
	sensorCache *CachedSensorStore
}

// Sensors returns a map of sensors
func (a *API) Sensors() (*Sensors, error) {

	url := fmt.Sprintf("%s/%s/sensors", a.Config.Addr, a.Config.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s: %s", url, err)
	}

	defer resp.Body.Close()

	var sensors Sensors

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&sensors)
	if err != nil {
		return nil, fmt.Errorf("unable to decode deCONZ response: %s", err)
	}

	return &sensors, nil

}

// EventReader returns a event.Reader with a default cached type store
func (a *API) EventReader() (*event.Reader, error) {

	if a.sensorCache == nil {
		a.sensorCache = &CachedSensorStore{SensorGetter: a}
	}

	if a.Config.wsAddr == "" {
		err := a.Config.discoverWebsocket()
		if err != nil {
			return nil, err
		}
	}

	return &event.Reader{TypeStore: a.sensorCache, WebsocketAddr: a.Config.wsAddr}, nil
}

// SensorEventReader takes an event reader and returns an sensor event reader
func (a *API) SensorEventReader(r *event.Reader) *SensorEventReader {

	if a.sensorCache == nil {
		a.sensorCache = &CachedSensorStore{SensorGetter: a}
	}

	return &SensorEventReader{lookup: a.sensorCache, reader: r}
}
