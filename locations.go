package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

type Location struct {
	UUID   string `db:"uuid"`
	UserID string `db:"userid"`
	Name   string `db:"name"`

	Lat  float64 `db:"lat"`
	Long float64 `db:"long"`
}

func (a api) CreateLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	decoder := json.NewDecoder(r.Body)
	location := Location{}

	err := decoder.Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	location.UserID = userID

	err = d.CreateLocation(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(location)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) ReadAllLocations(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	l, err := d.ReadAllLocations(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) ReadLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	locationID := ps.ByName("location")

	l, err := d.ReadLocation(userID, locationID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) UpdateLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	locationID := ps.ByName("location")

	decoder := json.NewDecoder(r.Body)
	location := Location{}

	err := decoder.Decode(&location)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	location.UserID = userID
	location.UUID = locationID

	err = d.UpdateLocation(location)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	return
}

func (a api) DeleteLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	locationID := ps.ByName("location")

	err := d.DeleteLocation(userID, locationID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}
}

func (d database) CreateLocation(l *Location) (err error) {
	l.UUID = uuid.NewV4().String()

	return d.InsertWithFields(l, "locations", "uuid", "userid", "name", "long", "lat")
}

func (d database) ReadAllLocations(userID string) (locations []Location, err error) {
	err = d.db.Select(&locations, "SELECT * FROM locations WHERE userid = $1",
		userID,
	)

	return
}

func (d database) ReadLocation(userID, locationID string) (l Location, err error) {
	locations, err := d.ReadAllLocations(userID)
	if err != nil {
		return
	}

	for _, l = range locations {
		if l.UUID == locationID {
			return
		}
	}

	err = fmt.Errorf("Returned %d locations, expected 1", len(locations))

	return
}

func (d database) UpdateLocation(l Location) (err error) {
	return d.UpdateWithFields(&l, "locations", "uuid", "userid", "name", "long", "lat")
}

func (d database) DeleteLocation(u, l string) (err error) {
	_, err = d.db.Exec("DELETE FROM locations WHERE userid = $1 AND uuid = $2", u, l)
	return
}
