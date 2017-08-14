package main

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
	_ "github.com/satori/go.uuid"
)

type database struct {
	db *sqlx.DB
}

func NewDatabase(connection string) (d database, err error) {
	d.db, err = sqlx.Connect("postgres", connection)

	return
}

func (d database) GetNearbyLocations(userID string, long, lat float64) (l []Location, err error) {
	err = d.db.Select(&l, "SELECT * FROM locations WHERE userid = $1 AND long > $2 AND long < $3 AND lat > $4 AND lat < $5",
		userID,
		long-0.02,
		long+0.02,
		lat-0.02,
		lat+0.02,
	)

	return
}

func (d database) GetTasks(userID, locationID string) (t []Task, err error) {
	err = d.db.Select(&t, "SELECT * FROM tasks WHERE userid = $1 AND locationid = $2",
		userID,
		locationID,
	)

	return
}

func (d database) InsertWithFields(i interface{}, table string, f ...string) (err error) {
	p := placeholders(f)

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(f, ", "),
		strings.Join(p, ", "),
	)

	_, err = d.db.NamedExec(query, i)

	return
}

func (d database) UpdateWithFields(i interface{}, table string, f ...string) (err error) {
	p := placeholders(f)

	query := fmt.Sprintf("UPDATE %s SET (%s) = (%s) WHERE uuid = :uuid",
		table,
		strings.Join(f, ", "),
		strings.Join(p, ", "),
	)

	_, err = d.db.NamedExec(query, i)

	return
}

func placeholders(f []string) (p []string) {
	for _, field := range f {
		p = append(p, fmt.Sprintf(":%s", field))
	}

	return
}
