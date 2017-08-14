package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a api) RegisterStream(userID string) (s Subscriber, err error) {
	_, err = d.ReadUser(userID)
	if err != nil {
		return
	}

	s = Subscriber{
		AMQP: NewAMQP(userID),
	}

	a.streamersSubscribers[userID] = s

	return
}

func (a api) Stream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var taskReceiver Subscriber
	var err error

	userID := ps.ByName("user")

	respf, ok := w.(http.Flusher)
	if !ok {
		panic("not flushable")
	}

	taskReceiver, ok = a.streamersSubscribers[userID]
	if !ok {
		taskReceiver, err = a.RegisterStream(userID)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	taskStream, err := taskReceiver.Receive(userID)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	for t := range taskStream {
		b := t.Body

		b = append(b, '\n')

		w.Write(b)
		respf.Flush()

		taskReceiver.ch.Ack(t.DeliveryTag, false)
	}
}
