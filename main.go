package main

import (
	"flag"
	"io/ioutil"
)

var (
	a api
	d database

	LoadSchema = flag.Bool("load-schema", false, "Whether to reload the schema into the db")
)

func main() {
	flag.Parse()

	var err error

	d, err = NewDatabase("postgres://ula:ula@database/tasks?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if *LoadSchema {
		schemaData, err := ioutil.ReadFile("db/schema.sql")
		if err != nil {
			panic(err)
		}

		d.db.MustExec(string(schemaData))
	}

	a, err = NewAPI(":8008")
	if err != nil {
		panic(err)
	}

	a.Start()
}
