package sensor

import (
	"time"
)

// Reading from the sensors for this application
type SensorReading struct {

	// The time this reading was taken
	Date time.Time

	// The moisture level of this reading
	Moisture float64
}

// Interface for storing sensor readings
type SensorStorer interface {
	
	// Stores the provided readings
	StoreReadings(readings []*SensorReading) bool

	// Gets the specified readings in the specified timeframe
	GetReadings(start time.Time, end time.Time) []*SensorReading

}