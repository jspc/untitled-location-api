package main

import (
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

	a.router.GET("/stream/:user/stream", a.Stream)

	a.router.PUT("/user/", a.CreateUser)
	a.router.GET("/user/:user", a.ReadUser)
	a.router.POST("/user/:user", a.UpdateUser)
	a.router.DELETE("/user/:user", a.DeleteUser)
	a.router.POST("/user/:user/checkin", a.CheckIn)

	a.router.PUT("/user/:user/location", a.CreateLocation)
	a.router.GET("/user/:user/location", a.ReadAllLocations)
	a.router.GET("/user/:user/location/:location", a.ReadLocation)
	a.router.POST("/user/:user/location/:location", a.UpdateLocation)
	a.router.DELETE("/user/:user/location/:location", a.DeleteLocation)

	a.router.PUT("/user/:user/task", a.CreateTask)
	a.router.GET("/user/:user/task", a.ReadAllTasks)
	a.router.GET("/user/:user/task/:task", a.ReadTask)
	a.router.POST("/user/:user/task/:task", a.UpdateTask)
	a.router.DELETE("/user/:user/task/:task", a.DeleteTask)

	return
}

func (a api) Start() {
	log.Panic(http.ListenAndServe(a.listen, a.router))
}
