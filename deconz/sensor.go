package deconz

// Sensors is a map of sensors indexed by their id
type Sensors map[int]Sensor

// Sensor is a deCONZ sensor, not that we only implement fields needed
// for event parsing to work
type Sensor struct {
	Type string
}
