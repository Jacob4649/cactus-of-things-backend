package sensor

import "math"

// Gets the mean and standard deviation of a set of sensor readings
func getMeanSD(readings []*SensorReading) (float64, float64) {
	mean := 0.0
	variance := 0.0
	length := len(readings)

	for _, reading := range readings {
		mean += reading.Moisture / float64(length)
	}

	for _, reading := range readings {
		variance += math.Pow(reading.Moisture - mean, 2) / float64(length)
	}

	return mean, math.Sqrt(variance)
}

// Assumes that readings are normally distributed around a sample value, and gets the z score corresponding to a reading
func getReadingZScore(reading SensorReading, mean float64, sd float64) float64 {
	if (sd == 0) {
		if (reading.Moisture == mean) {
			return 0
		} else {
			return 8192 / 100 // arbitrary, huge z score
		}
	}
	return (reading.Moisture - mean) / sd
}

// Determines if a single reading is from the same population as a group of sensor readings
// assumes the sensor readings are normally distributed
// uses a z score of 1.96 to represent the alpha significance of 0.05 for a 2 sided p test
func fromSamePopulation(reading SensorReading, readings []*SensorReading) bool {
	mean, sd := getMeanSD(readings)
	return math.Abs(getReadingZScore(reading, mean, sd)) < 1.96
}
