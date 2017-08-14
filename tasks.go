package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
)

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

func (a api) CreateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	decoder := json.NewDecoder(r.Body)
	task := Task{}

	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	task.UserID = userID

	err = d.CreateTask(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) ReadAllTasks(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")

	l, err := d.ReadAllTasks(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) ReadTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	taskID := ps.ByName("task")

	l, err := d.ReadTask(userID, taskID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	lB, err := json.Marshal(l)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	fmt.Fprintf(w, string(lB))
}

func (a api) UpdateTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	taskID := ps.ByName("task")

	decoder := json.NewDecoder(r.Body)
	task := Task{}

	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, "Invalid data")

		return
	}

	task.UserID = userID
	task.UUID = taskID

	err = d.UpdateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}

	return
}

func (a api) DeleteTask(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userID := ps.ByName("user")
	taskID := ps.ByName("task")

	err := d.DeleteTask(userID, taskID)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprintf(w, err.Error())

		return
	}
}

func (d database) CreateTask(l *Task) (err error) {
	l.UUID = uuid.NewV4().String()

	return d.InsertWithFields(l, "tasks", "uuid", "userid", "locationid", "type", "title", "description", "time")
}

func (d database) ReadAllTasks(userID string) (tasks []Task, err error) {
	err = d.db.Select(&tasks, "SELECT * FROM tasks WHERE userid = $1",
		userID,
	)

	return
}

func (d database) ReadTask(userID, taskID string) (l Task, err error) {
	tasks, err := d.ReadAllTasks(userID)
	if err != nil {
		return
	}

	for _, l = range tasks {
		if l.UUID == taskID {
			return
		}
	}

	err = fmt.Errorf("Returned %d tasks, expected 1", len(tasks))

	return
}

func (d database) UpdateTask(l Task) (err error) {
	return d.UpdateWithFields(&l, "tasks", "uuid", "userid", "locationid", "type", "title", "description", "time")
}

func (d database) DeleteTask(u, l string) (err error) {
	_, err = d.db.Exec("DELETE FROM tasks WHERE userid = $1 AND uuid = $2", u, l)
	return
}
