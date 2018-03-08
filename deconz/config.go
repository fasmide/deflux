package deconz

import (
	"net/url"
)

// Config represents a Deconz gateway
type Config struct {
	Addr   string
	APIKey string
	wsAddr url.URL
}

func (d *Config) discoverWebsocket() error {

	return nil
}
