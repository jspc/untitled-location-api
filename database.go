package main

import (
	"database/sql"

	_ "github.com/lib/pq"
	_ "github.com/satori/go.uuid"
)

type database struct {
	db *sql.DB
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
	d.db, err = sql.Open("postgres", connection)

	return
}

func (d database) GetNearbyLocations(userID string, long, lat float64) (l []Location, err error) {
	rows, err := d.db.Query("SELECT * FROM locations WHERE userid = $1 AND long > $2 AND long < $3 AND lat > $4 AND lat < $5",
		userID,
		long-0.001,
		long+0.001,
		lat-0.001,
		lat+0.001,
	)

	if err != nil {
		return
	}

	defer rows.Close()
	if err = rows.Err(); err != nil {
		return
	}

	for rows.Next() {
		var loc Location

		err = rows.Scan(
			&loc.UUID,
			&loc.UserID,
			&loc.Name,
			&loc.Long,
			&loc.Lat,
		)
		if err != nil {
			return
		}

		l = append(l, loc)
	}

	return
}
