package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	supabase "github.com/lengzuo/supa"
)

// User represents the structure of the data to be returned as JSON.
type User struct {
	ID   int    `json:"id"`
	Name string `json:"username"`
	Age  int    `json:"age"`
}

var temp = map[int]User{}

// Initialize Supabase client.
func initSupabaseClient() (supa *supabase.Client, err error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables instead.")
	}

	apiKey := os.Getenv("API_KEY")
	projectRef := os.Getenv("PROJECT_REF")
	debug := os.Getenv("DEBUG") == "true"

	if apiKey == "" || projectRef == "" {
		return nil, fmt.Errorf("API_KEY or PROJECT_REF is not set")
	}

	conf := supabase.Config{
		ApiKey:     apiKey,
		ProjectRef: projectRef,
		Debug:      debug,
	}
	return supabase.New(conf)
}

func fetchUsersFromSupabase(client *supabase.Client) (map[int]User, error) {
	users := make(map[int]User)

	var results []struct {
		ID   int    `json:"id"`
		Name string `json:"username"`
		Age  int    `json:"age"`
	}
	ctx := context.Background()
	err := client.DB.From("users").Select("id, username, age").Execute(ctx, &results)
	if err != nil {
		return nil, err
	}

	for _, user := range results {
		users[user.ID] = User{
			ID:   user.ID,
			Name: user.Name,
			Age:  user.Age,
		}
		fmt.Println(user)
	}

	return users, nil
}

func userHandler(w http.ResponseWriter, r *http.Request) {

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
	user, exists := temp[id]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func addUserHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var newUser struct {
		Username string `json:"username"`
		Age      int    `json:"age"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newUser.Username == "" || newUser.Age <= 0 {
		http.Error(w, "Invalid username or age", http.StatusBadRequest)
		return
	}

	supaClient, err := initSupabaseClient()
	if err != nil {
		http.Error(w, "Failed to initialize Supabase client", http.StatusInternalServerError)
		return
	}

	type User struct {
		Name string `json:"username"`
		Age  int    `json:"age"`
	}

	user := User{
		Name: newUser.Username,
		Age:  newUser.Age,
	}

	ctx := context.Background()
	err = supaClient.DB.From("users").Insert(user).Execute(ctx, nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s with age %d added successfully", newUser.Username, newUser.Age)
}

func main() {
	supaClient, err := initSupabaseClient()
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize Supabase client: %v", err))
	}

	users, err := fetchUsersFromSupabase(supaClient)
	temp = users
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch users from Supabase: %v", err))
	}

	http.HandleFunc("/username/", userHandler)
	http.HandleFunc("/adduser", addUserHandler)

	address := "localhost:8080"
	println("Server is running on", address)
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
