package main

import (
	"log"
	"os"
	"todo_app/app"
	"todo_app/app/auth"
	"todo_app/app/todo-app"
	"todo_app/common"
	"todo_app/config"
	"todo_app/store"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
		os.Exit(1)
	}
	// setup logger
	logger, err := common.NewLogger(conf.Env)
	defer logger.Sync()
	if err != nil {
		log.Fatalf("logging: %v", err)
		os.Exit(1)
	}
	logger.Info("starting the application...")

	// setup pg store
	store, err := store.New(conf.Store.ConnString)
	if err != nil {
		logger.Fatalw("couldn't instantiate store", "reason", err.Error())
		os.Exit(1)
	}

	// run migrations
	msg, err := store.MigrateDatabase(conf.Store.MigrationsPath)
	if err != nil {
		logger.Fatalw("couldn't run migrations", "reason", err.Error())
		os.Exit(1)
	}
	logger.Info(msg)

	// setup app
	app := app.New()

	app.Use(common.RequestLoggingMiddleware(logger))

	auth.UseEndpoints(app, store, conf.Application.Secret)

	app.Use(auth.JwtAuthMiddleware()) // after this middleware, every request will be checked for jwt

	todo.UseEndpoints(app, store)

	if err = app.Listen(conf.Application.HostUrl); err != nil {
		logger.Fatalw("application failed", "reason", err.Error())
	}
}
