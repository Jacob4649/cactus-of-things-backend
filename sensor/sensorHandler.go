package sensor

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Handler function
type Handler func(http.ResponseWriter, *http.Request)

// Gets handlers to read sensor values and write sensor values with this sensor storer
func GetHandlers(storer SensorStorer) (getter Handler, setter Handler) {

	getHandler := func (w http.ResponseWriter, r *http.Request)  {
		
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

		successs := storer.StoreReadings(readings)

		if !successs {
			msg := "Unable to write to database"
			http.Error(w, msg, http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, "Readings written")	
		}
	}

	return getHandler, setHandler
}