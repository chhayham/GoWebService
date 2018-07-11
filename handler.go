package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Index GET home page
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

// TodoIndex GET todo list
func TodoIndex(w http.ResponseWriter, r *http.Request) {
	todos := Todos{
		Todo{Name: "Write presentation"},
		Todo{Name: "Host meetup"},
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		panic(err)
	}
}

// TodoShow GET todo entry
func TodoShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoID)
}

// CreateTodo POST new todo entry
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Created task:")
}

// UpdateTodo PUT update todo entry
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Updated task:")
}

// DeleteTodo DELETE delete todo entry
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Deleted task:")
}
