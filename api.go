package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

type api struct {
	router *httprouter.Router
	listen string

	streamersTasks map[string]map[string]chan Task
}

func NewAPI(listen string) (a api, err error) {
	a.listen = listen
	a.streamersTasks = make(map[string]map[string]chan Task)
	a.router = httprouter.New()

	a.router.POST("/:user/checkin", a.CheckIn)

	a.router.POST("/:user/stream", a.RegisterStream)
	a.router.GET("/:user/stream/:streamID", a.Stream)

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
				for _, stream := range a.streamersTasks[userID] {
					stream <- task
				}
			}
		}
	}
}

func (a api) RegisterStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	_, ok := a.streamersTasks[userID]
	if !ok {
		a.streamersTasks[userID] = map[string]chan Task{}
	}

	streamID := uuid.NewV4().String()
	a.streamersTasks[userID][streamID] = make(chan Task)

	fmt.Fprintf(w, streamID)
}

func (a api) Stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	streamID := ps.ByName("streamID")

	respf, ok := w.(http.Flusher)
	if !ok {
		panic("not flushable")
	}

	taskStream, ok := a.streamersTasks[userID][streamID]
	if !ok {
		w.WriteHeader(404)
		return
	}

	defer delete(a.streamersTasks[userID], streamID)

	w.WriteHeader(200)
	for t := range taskStream {
		b, err := json.Marshal(t)
		if err != nil {
			panic(err)
		}

		b = append(b, '\n')

		w.Write(b)
		respf.Flush()
	}
}
