package deconz

import (
	"fmt"
	"time"

	"github.com/fasmide/deflux/deconz/event"
)

// SensorLookup represents an interface for sensor lookup
type SensorLookup interface {
	LookupSensor(int) (*Sensor, error)
}

// EventReader interface
type EventReader interface {
	ReadEvent() (*event.Event, error)
}

// SensorEventReader reads events from an event.reader and returns SensorEvents
type SensorEventReader struct {
	lookup SensorLookup
	reader EventReader
	running bool
}






// starts a thread reading events into the given channel
// returns immediately
func (r *SensorEventReader) Start(out chan *SensorEvent) error {
	
	if r.lookup == nil {
		return errors.New("Cannot run without a SensorLookup from which to lookup sensors")
	}
	if r.reader == nil {
		return errors.New("Cannot run without a EventReader from which to read events")
	}

	if r.running {
		return "Reader is already running."
	}

	r.running = true
	
	go func() {
		REDIAL:
		for r.running {
			// establish connection
			for r.running {
				err := r.reader.Dial()
				if err != nil {
					log.Printf("Error connecting Deconz websocket: %s\nAttempting reconnect in 5s...", err)
					time.Sleep(5 * time.Second) // TODO configurable delay
				} else {
					log.Printf("Deconz websocket connected")
				}
			}
			// read events until connection fails
			for r.running {
				e, err := r.reader.ReadEvent()
				if err != nil {
					if eerr, ok := err.(EventError) ; ok && eerr.Recoverable() {
						log.Printf("Dropping event due to error: %s", err)
						continue
					}
					continue REDIAL
				}

				sensor, err := s.lookup.LookupSensor(e.ID)
				if err != nil {
					log.Printf("Dropping event. Could not lookup sensor for id %d: %s", e.ID, err)
					continue
				}
				// send event on channel
				out <- &SensorEvent{Event: e, Sensor: sensor}
			}
		}
		// if not running, close connection and return from goroutine
		r.reader.Close()
		log.Printf("Deconz websocket closed")
	}()
	return nil
}



// Close closes the reader, closing the connection to deconz and terminating the goroutine
func (r *SensorEventReader) StopReadEvents() {
	r.running = false
}
