package sensor

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
)

// Datastore implementation of SensorStorer
type SensorDatastore struct {

	// ID of GCP project
	ProjectID string

	// Client for the project
	client *datastore.Client

}

// Datastore sensor reading
type SensorDatastoreReading struct {
	
	// The time this reading was taken
	Date time.Time `datastore:"date"`
	
	// The moisture level of this reading
	Moisture float64 `datastore:"moisture,noindex"`

	// The light level of this reading
	Light float64 `datastore:"light,noindex"`

	// The expiry date for this reading
	Expiry time.Time `datastore:"expiry,noindex"`
}

// Sets up client for sensor datastore
func (store *SensorDatastore) initializeClient() (*datastore.Client, error) {
	context := context.Background()

	var client *datastore.Client = store.client
	var err error = nil

	if client == nil {
		client, err = datastore.NewClient(context, store.ProjectID)
	}

	return client, err
}

func (store *SensorDatastore) GetReadings(start time.Time, end time.Time) []*SensorReading {
	if start == end || end.Before(start) {
		return []*SensorReading{} // early return on invalid intervals
	}

	context := context.Background()

	client, err := store.initializeClient()

	if err != nil {
		return nil
	}

	kind := "cactus-sensor-reading"

	query := datastore.NewQuery(kind).FilterField("date", ">=", start).FilterField("date", "<=", end)
	
	datastoreReadings := []*SensorDatastoreReading{}

	_, err = client.GetAll(context, query, &datastoreReadings)

	if err != nil {
		return nil
	}

	readings := []*SensorReading{}

	for _, reading := range datastoreReadings {
		readings = append(readings, &SensorReading{Date: reading.Date, Moisture: reading.Moisture})
	}

	return readings
}

func (store *SensorDatastore) StoreReadings(readings []*SensorReading) bool {
	
	if len(readings) == 0 {
		return true; // early return on empty input
	}

	context := context.Background()

	client, err := store.initializeClient()

	if err != nil {
		return false
	}
	
	defer client.Close() // close the client after this function

	kind := "cactus-sensor-reading"

	keys := []*datastore.Key{}

	datastoreReadings := []*SensorDatastoreReading{}

	for _, reading := range readings {
		keys = append(keys, datastore.IncompleteKey(kind, nil))
		datastoreReadings = append(datastoreReadings, &SensorDatastoreReading{Moisture: reading.Moisture,
			Light: reading.Light,
			Date: reading.Date,
			Expiry: reading.Date.Add(time.Hour * 24 * 30)})
	}

	if keys, err = client.AllocateIDs(context, keys); err != nil {
		return false
	}

	keys, err = client.PutMulti(context, keys, datastoreReadings);

	return err == nil
}