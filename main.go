package main

var (
	a api
	d database
)

func main() {
	var err error

	d, err = NewDatabase("postgres://postgres@database/tasks?sslmode=disable")
	if err != nil {
		panic(err)
	}

	a, err = NewAPI(":8008")
	if err != nil {
		panic(err)
	}

	a.Start()
}
