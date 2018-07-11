package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route struct for url route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes slice of route containing name method path and handler.  see handler.go
type Routes []Route

// NewRouter gorilla mux
func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

// CRUD REST API routes
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"CreateTodo",
		"POST",
		"/create",
		CreateTodo,
	},
	Route{
		"UpdateTodo",
		"PUT",
		"/update",
		UpdateTodo,
	},
	Route{
		"DeleteTodo",
		"DELETE",
		"/delete",
		DeleteTodo,
	},
}
