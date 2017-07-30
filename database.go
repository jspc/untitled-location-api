package main

import (
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "github.com/satori/go.uuid"
)

type database struct {
	db *sqlx.DB
}

type User struct {
	UUID  string
	Email string
}

type Location struct {
	UUID   string
	UserID string
	Name   string

	Long float64
	Lat  float64
}

type Task struct {
	UUID        string
	UserID      string
	LocationID  string
	Type        int
	Title       string
	Description string
	Time        string
}

// Task Types
const (
	SimpleTask = iota
)

func NewDatabase(connection string) (d database, err error) {
	d.db, err = sqlx.Connect("postgres", connection)

	return
}

func (d database) GetNearbyLocations(userID string, long, lat float64) (l []Location, err error) {
	err = d.db.Select(&l, "SELECT * FROM locations WHERE userid = $1 AND long > $2 AND long < $3 AND lat > $4 AND lat < $5",
		userID,
		long-0.001,
		long+0.001,
		lat-0.001,
		lat+0.001,
	)

	return
}
