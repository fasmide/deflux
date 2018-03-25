package deconz

import (
	"errors"
	"fmt"
)

// CachedTypeStore is a cached typestore which provides LookupType for event passing
// it will be our default store
type CachedTypeStore struct {
	SensorGetter
	cache *Sensors
}

// SensorGetter defines how we like to ask for sensors
type SensorGetter interface {
	Sensors() (*Sensors, error)
}

// LookupType lookups deCONZ event types though a cache
// TODO: if we where unable to lookup an ID we should try to refetch the cache
// - there could have been an sensor added we dont know about
func (c *CachedTypeStore) LookupType(i int) (string, error) {
	var err error
	if c.cache == nil {
		err = c.populateCache()
		if err != nil {
			return "", fmt.Errorf("unable to populate types: %s", err)
		}
	}

	if s, found := (*c.cache)[i]; found {
		return s.Type, nil
	}

	return "", errors.New("no such sensor")
}

func (c *CachedTypeStore) populateCache() error {
	var err error
	c.cache, err = c.Sensors()
	if err != nil {
		return err
	}
	return nil
}
