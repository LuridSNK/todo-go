package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

var (
	address string
)

var storage = map[uuid.UUID]todoItem{}

func main() {
	address = "localhost:5000"
	logger := log.New(os.Stdout, "http-server: ", log.LstdFlags)

	server := &http.Server{
		Addr:         address,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      logRequestMiddleware(logger)(defineRoutes()),
	}

	logger.Printf("Listening at %s", address)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Couldn't listen on %s %v\n", address, err)
	}
}

func defineRoutes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/todo", routeToDoItems)
	return router
}

func routeToDoItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodoItems(w, r)
		return
	case "POST":
		addNewTodoItem(w, r)
		return
	}
}

func getTodoItems(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	keys := make([]todoItem, 0, len(storage))
	for k := range storage {
		keys = append(keys, storage[k])
	}
	response, err := json.Marshal(keys)
	if err != nil {
		fmt.Println(w, err)
		return
	}

	fmt.Fprintf(w, "%s", response)
}

func addNewTodoItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var i todoItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&i)
	i.Id = uuid.New()
	i.CreatedAt = time.Now().UTC()
	if err != nil {
		fmt.Println(w, err)
		fmt.Fprintf(w, "Could not create an item, reason: %s", err)
	}

	storage[i.Id] = i

	fmt.Fprintln(w, "CREATED")
}

func removeTodoItem(w http.ResponseWriter, r *http.Request) {
	// remove item is not possible to do with standard net/http without making an overhead wrapper
}

func logRequestMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				logger.Println(uuid.New(), r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			}()
			next.ServeHTTP(w, r)
		})
	}
}

type todoItem struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	IsDone      bool      `json:"isDone"`
}
