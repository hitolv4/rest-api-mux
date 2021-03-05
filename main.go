package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}
type AllTasks []Task

var tasks = AllTasks{
	{
		ID:      1,
		Name:    "Task One",
		Content: "Some Content",
	},
	{
		ID:      2,
		Name:    "Task Two",
		Content: "Some content",
	},
	{
		ID:      3,
		Name:    "Task Two",
		Content: "Some content",
	},
	{
		ID:      4,
		Name:    "Task Two",
		Content: "Some content",
	},
	{
		ID:      5,
		Name:    "Task Two",
		Content: "Some content",
	},
}

func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API!!!")
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}
func GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(w, r)
	if err != nil {
		fmt.Fprint(w, "id invalid")
		return
	}
	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusFound)
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "not found")
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(w, r)
	if err != nil {
		fmt.Fprint(w, "id invalid")
		return
	}
	for i, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-type", "application/json")
			//w.WriteHeader(http.StatusNoContent)
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with ID %v has been removed succesfully", id)
			return

		}
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "not found")

}

func GetID(w http.ResponseWriter, r *http.Request) (int, error) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := GetID(w, r)
	if err != nil {
		fmt.Fprint(w, "id invalid")
		return
	}
	var updatetask Task
	requestBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(requestBody, &updatetask)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task")
	}
	for i, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			updatetask.ID = id
			if updatetask.Name == "" {
				updatetask.Name = task.Name
			}
			if updatetask.Content == "" {
				updatetask.Content = task.Content
			}
			tasks[i] = updatetask
			fmt.Fprintf(w, "The task with ID %v has been updated succesfully", id)
			return

		}
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Insert a Valid Task")
	}
	json.Unmarshal(requestBody, &newTask)
	newTask.ID = tasks[len(tasks)-1].ID + 1
	tasks = append(tasks, newTask)

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func main() {
	fmt.Println("Server Started")
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexRoute)
	r.HandleFunc("/tasks", GetTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", GetTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")
	r.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks", CreateTask).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", r))
}
