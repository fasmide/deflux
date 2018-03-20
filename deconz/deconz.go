package deconz

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Deconz represents a deconz server device
type Deconz struct {
	sync.Mutex
	Config   Config
	handlers []Handler
}

// Handler is the required interface for handlers
type Handler func(interface{})

// Connect connects and calls all handlers for all handlers
func (d *Deconz) Connect() error {

	c, _, err := websocket.DefaultDialer.Dial(d.Config.wsAddr, nil)
	if err != nil {
		return fmt.Errorf("unable to dail websocket: %s", err)
	}
	log.Printf("Connected!")
	defer c.Close()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %s", err)
		}
		log.Printf("recv: %s", message)
	}
}

// AddHandler adds a new handler that should receive events
func (d *Deconz) AddHandler(h Handler) {
	d.Lock()
	d.handlers = append(d.handlers, h)
	d.Unlock()
}
