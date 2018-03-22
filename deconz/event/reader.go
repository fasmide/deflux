package event

import (
	"errors"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

// Reader represents a deconz server device
type Reader struct {
	WebsocketAddr string
	TypeStore     TypeLookuper
	decoder       *Decoder
	conn          *websocket.Conn
}

// Dial connects connects to deconz, use ReadEvent to recieve events
func (r *Reader) Dial() error {

	// // if wsAddr is empty, discover it..
	// if r.Config.wsAddr == "" {
	// 	err := r.Config.discoverWebsocket()
	// 	if err != nil {
	// 		return fmt.Errorf("unable to dail websocket: %s", err)
	// 	}
	// }

	if r.TypeStore == nil {
		return errors.New("cannot dial without a TypeStore to lookup events from")
	}

	// create a decoder with the typestore
	r.decoder = &Decoder{TypeStore: r.TypeStore}

	// connect
	var err error
	r.conn, _, err = websocket.DefaultDialer.Dial(r.WebsocketAddr, nil)
	if err != nil {
		return fmt.Errorf("unable to dail websocket: %s", err)
	}
	return nil
}

// ReadEvent reads, parses and returns the next event
func (r *Reader) ReadEvent() (*Event, error) {

	_, message, err := r.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("event read error: %s", err)
	}

	log.Printf("recv: %s", message)

	e, err := r.decoder.Parse(message)
	if err != nil {
		return nil, fmt.Errorf("unable to parse message: %s", err)
	}

	return e, nil
}

// Close closes the connection to deconz
func (r *Reader) Close() error {
	return r.conn.Close()
}
