package deconz

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// DeconzDiscoveryEndpoint is the url used when auto discovering a deconz gateway
const DeconzDiscoveryEndpoint = "https://dresden-light.appspot.com/discover"

// DiscoveryResponse is a slice of discovered gateways
type DiscoveryResponse []Discovery

// Discovery is a discovered deconz gateway
// [{"macaddress": "00212EFFFF017FBD", "name": "deCONZ-GW", "internalipaddress": "192.168.1.90", "publicipaddress": "85.191.222.130", "internalport": 8080, "id": "00212EFFFF017FBD"}]
type Discovery struct {
	ID                string
	Name              string
	MacAddress        string
	PublicIPAddress   string
	InternalIPAddress string
	InternalPort      uint
}

// Discover discovers deconz gateways
func Discover() (DiscoveryResponse, error) {
	response, err := http.Get(DeconzDiscoveryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to talk to discovery endpoint: %s", err)
	}

	var data DiscoveryResponse

	d := json.NewDecoder(response.Body)
	defer response.Body.Close()

	err = d.Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("unable to parse json from discovery endpoint: %s", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no gateways was found")
	}

	return data, nil
}
