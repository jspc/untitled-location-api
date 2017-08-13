package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type api struct {
	router *httprouter.Router
	listen string

	streamersSubscribers map[string]Subscriber
}

type message struct {
	Type string
	Body interface{}
}

func NewAPI(listen string) (a api, err error) {
	a.listen = listen
	a.streamersSubscribers = make(map[string]Subscriber)

	a.router = httprouter.New()

	a.router.POST("/user/:user/checkin", a.CheckIn)

	a.router.POST("/user/:user/stream", a.RegisterStream)
	a.router.GET("/user/:user/stream", a.Stream)

	a.router.PUT("/user/", a.CreateUser)
	a.router.GET("/user/:user", a.ReadUser)
	a.router.POST("/user/:user", a.UpdateUser)
	a.router.DELETE("/user/:user", a.DeleteUser)

	a.router.PUT("/user/:user/location", a.CreateLocation)
	a.router.GET("/user/:user/location", a.ReadAllLocations)
	a.router.GET("/user/:user/location/:location", a.ReadLocation)
	a.router.POST("/user/:user/location/:location", a.UpdateLocation)
	a.router.DELETE("/user/:user/location/:location", a.DeleteLocation)

	return
}

func (a api) Start() {
	log.Panic(http.ListenAndServe(a.listen, a.router))
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

func (a api) RegisterStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	_, ok := a.streamersSubscribers[userID]
	if !ok {
		a.streamersSubscribers[userID] = Subscriber{
			AMQP: NewAMQP(userID),
		}
	}
}

func (a api) Stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	respf, ok := w.(http.Flusher)
	if !ok {
		panic("not flushable")
	}

	taskReceiver, ok := a.streamersSubscribers[userID]
	if !ok {
		w.WriteHeader(404)
		return
	}

	taskStream, err := taskReceiver.Receive(userID)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	for t := range taskStream {
		b := t.Body

		b = append(b, '\n')

		w.Write(b)
		respf.Flush()

		taskReceiver.ch.Ack(t.DeliveryTag, false)
	}
}
