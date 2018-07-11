package main

import "time"
// Todo model data
type Todo struct {
    Name      string    `json:"name"`
    Completed bool      `json:"completed"`
    Due       time.Time `json:"due"`
}
// Todos slice
type Todos []Todo