package main

import (
	"log"
	"os"
	"todo_app/app"
	"todo_app/app/auth"
	"todo_app/app/todo-app"
	"todo_app/common"
	"todo_app/store"
)

var (
	hostUrl          string
	storeUrl         string
	migrationsOutput string
	environment      string
)

func main() {
	// temporary setup env
	hostUrl = "localhost:5000"
	storeUrl = "postgresql://postgres:secret@localhost"
	migrationsOutput = "./store/migrations"
	environment = "development"

	// setup logger
	logger, err := common.NewLogger(environment)
	defer logger.Sync()
	if err != nil {
		log.Fatalf("logging: %v", err)
		os.Exit(1)
	}
	logger.Info("starting the application...")

	// setup pg store
	store, err := store.New(storeUrl)
	if err != nil {
		logger.Fatalw("couldn't instantiate store", "reason", err.Error())
		os.Exit(1)
	}

	// run migrations
	msg, err := store.MigrateDatabase(migrationsOutput)
	if err != nil {
		logger.Fatalw("couldn't run migrations", "reason", err.Error())
		os.Exit(1)
	}
	logger.Info(msg)

	// setup app
	app := app.New()

	app.Use(common.RequestLoggingMiddleware(logger))

	auth.UseEndpoints(app, store)

	app.Use(auth.JwtAuthMiddleware()) // after this middleware, every request will be checked for jwt

	todo.UseEndpoints(app, store)

	if err = app.Listen(hostUrl); err != nil {
		logger.Fatalw("application failed", "reason", err.Error())
	}
}
