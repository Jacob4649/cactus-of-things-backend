package sensor

import (
	"math"
	"time"
)

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

// Estimates the population standard deviation
// uses the pooled estimate of the standard deviation of the sampling distribution
func getPopulationSDEstimate(sd float64, n int) float64 {

	// n_1 = n
	// n_2 = 1
	// sd_1 = sd
	// sd_2 = 0

	degreesFreedom := n - 1 // n_1 + n_2 - 2
	
	sumN := n + 1 // n_1 + n_2

	// remembering that n_2s_2^2 = 0

	top := float64(n) * sd * sd * float64(sumN) // (n_1*sd_1^2 + n_2*sd_2^2) * (n_1 + n_2)

	bottom := degreesFreedom * n // degree of freedom * n_1 * n_2

	ratio := top / float64(bottom) // pooled estimate of sampling variance

	return math.Sqrt(ratio) // pooled estimate of standard deviation
}

// Assumes that readings are normally distributed around a sample value, and gets the t score corresponding to a reading
func getReadingTScore(reading SensorReading, mean float64, sd float64, n int) float64 {
	if (sd == 0) { // handles cases where you would be dividing by 0
		if (reading.Moisture == mean) {
			return 0
		} else {
			return 8192 / 100 // arbitrary, huge z score
		}
	}
	return (reading.Moisture - mean) / getPopulationSDEstimate(sd, n)
}

// gets the t score to use with the t distribution with the specified degrees of freedom
// for a two tailed test at 0.01 significance
// really shitty estimate, horrible for both very low and very large values
// most importantly, accurate-ish aroung 5, which will proabbly be the most common degree of freedom
func getTScore(df int) float64 {
	return 1 / math.Pow(float64(df), 2.3) * (63.6651 - 2.576) + 2.576
}

// Determines if a single reading is from the same population as a group of sensor readings
// assumes the sensor readings are normally distributed
// uses a approximate t score to represent the alpha significance of 0.05 for a 2 sided p test
func WateringFromSamePopulation(reading SensorReading, readings []*SensorReading) bool {
	mean, sd := getMeanSD(readings)

	length := len(readings)

	score := getReadingTScore(reading, mean, sd, length)

	return math.Abs(score) < getTScore(length - 1)
}

// Interface for storing information related to watering
type WaterStorer interface {

	// Stores that the plant was watered at a specific time
	StoreWatering(time.Time)

}
