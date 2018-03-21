package deconz

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

// Config represents a Deconz gateway
type Config struct {
	Addr   string
	APIKey string
	wsAddr string
}

// config is used to parse the things we need from the deCONZ config endpoint
type config struct {
	Websocketport int
}

func (c *Config) discoverWebsocket() error {
	u, err := url.Parse(c.Addr)
	if err != nil {
		return fmt.Errorf("unable to discover websocket: %s", err)
	}
	u.Path = path.Join(u.Path, c.APIKey, "config")

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("unable to discover websocket: %s", err)
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)

	var conf config
	err = dec.Decode(&conf)
	if err != nil {
		return fmt.Errorf("unable to discover websocket: %s", err)
	}

	// change our old parsed url to websocket, it should connect to the websocket endpoint of deCONZ
	u.Scheme = "ws"
	u.Path = "/"
	u.Host = fmt.Sprintf("%s:%d", u.Hostname(), conf.Websocketport)

	c.wsAddr = u.String()
	return nil
}
