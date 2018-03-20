package deconz

// Config represents a Deconz gateway
type Config struct {
	Addr   string
	APIKey string
	wsAddr string
}

func (d *Config) discoverWebsocket() error {

	return nil
}
