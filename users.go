package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

type User struct {
	UUID  string `db:"uuid"`
	Email string `db:"email"`
}

func (a api) CreateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	user := User{}

	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	err = d.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	uB, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(uB))
}

func (a api) ReadUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	u, err := d.ReadUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	uB, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(uB))
}

func (a api) UpdateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	decoder := json.NewDecoder(r.Body)
	user := User{}

	err := decoder.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	user.UUID = userID

	err = d.UpdateUser(user)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	return
}

func (a api) DeleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	err := d.DeleteUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}
}

func (a api) CheckIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	decoder := json.NewDecoder(r.Body)
	longLat := make(map[string]float64)

	err := decoder.Decode(&longLat)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	long := longLat["long"]
	lat := longLat["lat"]

	l, err := d.GetNearbyLocations(userID, long, lat)
	if err != nil {
		panic(err)
	}

	for _, location := range l {
		if Distance(Location{Lat: lat, Long: long}, location) <= 50 {

			t, err := d.GetTasks(userID, location.UUID)
			if err != nil {
				panic(err)
			}

			for _, task := range t {
				Publisher{NewAMQP(userID)}.Publish(message{
					Type: "task-in-range",
					Body: task,
				})
			}
		}
	}
}

func (d database) CreateUser(u *User) (err error) {
	u.UUID = uuid.NewV4().String()

	return d.InsertWithFields(u, "users", "uuid", "email")
}

func (d database) ReadUser(uuid string) (u User, err error) {
	users := []User{}

	err = d.db.Select(&users, "SELECT * FROM users WHERE uuid = $1",
		uuid,
	)

	if err != nil {
		return
	}

	if len(users) != 1 {
		err = fmt.Errorf("Returned %d users, expected 1", len(uuid))

		return
	}
	u = users[0]

	return
}

func (d database) UpdateUser(u User) (err error) {
	return d.UpdateWithFields(&u, "users", "uuid", "email")
}

func (d database) DeleteUser(u string) (err error) {
	_, err = d.db.Exec("DELETE FROM users WHERE uuid = $1", u)
	return
}
