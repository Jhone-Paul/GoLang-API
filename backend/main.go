package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// User represents the structure of the data to be returned as JSON.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Sample data to simulate a database.
var users = map[int]User{
	1: {ID: 1, Name: "Alice", Age: 30},
	2: {ID: 2, Name: "Bob", Age: 25},
	3: {ID: 3, Name: "Charlie", Age: 35},
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL.
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid request URL", http.StatusBadRequest)
		return
	}

	idStr := pathParts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Fetch the user from the sample data.
	user, exists := users[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Respond with the user data in JSON format.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/username/", userHandler)

	// Start the server on localhost:8080.
	address := "localhost:8080"
	println("Server is running on", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
