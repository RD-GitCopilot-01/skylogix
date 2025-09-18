package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoStore struct {
	todos  []Todo
	nextID int
	mutex  sync.RWMutex
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos:  make([]Todo, 0),
		nextID: 1,
	}
}

func (ts *TodoStore) GetAll() []Todo {
	ts.mutex.RLock()
	defer ts.mutex.RUnlock()
	return ts.todos
}

func (ts *TodoStore) Create(title string) Todo {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	todo := Todo{
		ID:        ts.nextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	
	ts.todos = append(ts.todos, todo)
	ts.nextID++
	
	return todo
}

func (ts *TodoStore) Update(id int, completed bool) (*Todo, error) {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	for i, todo := range ts.todos {
		if todo.ID == id {
			ts.todos[i].Completed = completed
			return &ts.todos[i], nil
		}
	}
	
	return nil, fmt.Errorf("todo with id %d not found", id)
}

func (ts *TodoStore) Delete(id int) error {
	ts.mutex.Lock()
	defer ts.mutex.Unlock()
	
	for i, todo := range ts.todos {
		if todo.ID == id {
			ts.todos = append(ts.todos[:i], ts.todos[i+1:]...)
			return nil
		}
	}
	
	return fmt.Errorf("todo with id %d not found", id)
}

var store = NewTodoStore()

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todos := store.GetAll()
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var req struct {
		Title string `json:"title"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	
	todo := store.Create(req.Title)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	var req struct {
		Completed bool `json:"completed"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	todo, err := store.Update(id, req.Completed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	
	if err := store.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	
	r.HandleFunc("/api/todos", getTodos).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id}", updateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{id}", deleteTodo).Methods("DELETE")
	
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})
	
	handler := c.Handler(r)
	
	fmt.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
