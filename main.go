package main

import (
	"fmt"
	"log"
	"os"
	"todo_app/app"
	"todo_app/app/auth"
	"todo_app/app/todo-app"
	"todo_app/store"

	"github.com/gofiber/fiber/v2/middleware/logger"
)

var (
	hostUrl          string
	storeUrl         string
	migrationsOutput string
)

func main() {
	output := os.Stderr
	hostUrl = "localhost:5000"
	storeUrl = "postgresql://postgres:secret@localhost"
	migrationsOutput = "./store/migrations"
	app := app.New()
	store, err := store.New(storeUrl)
	if err != nil {
		fmt.Fprintf(output, "FATAL [store]: %v", err)
		os.Exit(1)
	}

	err = store.MigrateDatabase(migrationsOutput)
	if err != nil {
		fmt.Fprintf(output, "FATAL [migrate]: %v", err)
		os.Exit(1)
	}

	app.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02T15:04:05",
		TimeZone:   "UTC",
		Output:     output,
		Format:     "[${time}] ${status} - ${latency} ${method} ${path} ${body}",
	}))

	todo.UseEndpoints(app, store)
	auth.UseEndpoints(app, store)

	log.Fatal(app.Listen(hostUrl))
}
