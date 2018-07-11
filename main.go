/*
Go REST api example.  Ref https://thenewstack.io/make-a-restful-json-api-go/
*/

package main

import (
	"log"
	"net/http"
)

func main() {
	router := NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
