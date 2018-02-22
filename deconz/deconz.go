package deconz

import (
	"fmt"
	"net/url"
)

// Config represents a Deconz gateway
type Config struct {
	Addr   string
	APIKey string
	wsAddr url.URL
}

// Init initializes endpoints
func (d *Config) Init() error {

	// in any case, talk to the api about where the websocket is located at...
	err := d.discoverWebsocket()
	if err != nil {
		return fmt.Errorf("unable to discover websocket endpoint: %s", err)
	}

	return nil
}

func (d *Config) discoverWebsocket() error {
	return nil
}

func (d *Config) discover() error {
	return nil
}
