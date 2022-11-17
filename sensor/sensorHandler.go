package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Handler function
type Handler func(http.ResponseWriter, *http.Request)

// Gets handlers to read sensor values and write sensor values with this sensor storer
func GetHandlers(storer SensorStorer) (getter Handler, setter Handler) {

	getHandler := func (w http.ResponseWriter, r *http.Request)  {
		parameters := r.URL.Query()
		
		startString, endString := parameters.Get("start"), parameters.Get("end")

		if startString == "" || endString == "" {
			msg := "Must specify start and end query strings in URL parameters"
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		start, startErr := time.Parse(time.RFC3339, startString)

		end, endErr := time.Parse(time.RFC3339, endString)

		if startErr != nil || endErr != nil {
			msg := "Start and end query parameters must be according to RFC3339"
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		readings := storer.GetReadings(start, end)

		if readings == nil {
			msg := "Unable to read from database"
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Add("Access-Control-Allow-Origin", "*");

		json.NewEncoder(w).Encode(readings)
	}

	setHandler := func (w http.ResponseWriter, r *http.Request)  {
		
		if r.Header.Get("Content-Type") != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}

		var readings []*SensorReading

		err := json.NewDecoder(r.Body).Decode(&readings)

		if err != nil {
			msg := "Unable to parse request json"
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		setDates(&readings)

		success := storer.StoreReadings(readings)

		if !success {
			msg := "Unable to write to database"
			http.Error(w, msg, http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, "Readings written")	
		}
	}

	return getHandler, setHandler
}

func setDates(readings *[]*SensorReading) {
	for _, reading := range *readings {
		if reading.Date.IsZero() {
			reading.Date = time.Now()
		}
	}
}
