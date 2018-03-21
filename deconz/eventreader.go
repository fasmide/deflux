package deconz

import (
	"fmt"
	"log"

	"github.com/fasmide/deflux/deconz/event"
	"github.com/gorilla/websocket"
)

// EventReader represents a deconz server device
type EventReader struct {
	Config Config
	conn   *websocket.Conn
}

// Dial connects connects to deconz, use ReadEvent to recieve events
func (d *EventReader) Dial() error {

	// if wsAddr is empty, discover it..
	if d.Config.wsAddr == "" {
		err := d.Config.discoverWebsocket()
		if err != nil {
			return fmt.Errorf("unable to dail websocket: %s", err)
		}
	}

	// connect
	var err error
	d.conn, _, err = websocket.DefaultDialer.Dial(d.Config.wsAddr, nil)
	if err != nil {
		return fmt.Errorf("unable to dail websocket: %s", err)
	}
	return nil
}

// ReadEvent reads, parses and returns the next event
func (d *EventReader) ReadEvent() (*event.Event, error) {
	_, message, err := d.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("event read error: %s", err)
	}

	log.Printf("recv: %s", message)

	e, err := event.Parse(message)
	if err != nil {
		return nil, fmt.Errorf("unable to parse message: %s", err)
	}

	return e, nil
}

// Close closes the connection to deconz
func (d *EventReader) Close() error {
	return d.conn.Close()
}
