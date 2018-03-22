package deconz

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// API represents the deCONZ rest api
type API struct {
	Config Config
}

// GetSensors returns a map of sensors
func (a *API) GetSensors() (*Sensors, error) {
	url := fmt.Sprint("%s/%s/sensors", a.Config.Addr, a.Config.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to get %s: %s", url, err)
	}

	defer resp.Body.Close()

	var sensors Sensors

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&sensors)
	if err != nil {
		return nil, fmt.Errorf("unable to parse deCONZ response: %s", err)
	}
	return

}
