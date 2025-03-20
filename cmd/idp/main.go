package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/guillaumebour/ghostidp/internal/application"
	"github.com/guillaumebour/ghostidp/internal/handlers"
	"github.com/guillaumebour/ghostidp/internal/utils/logger"
	"github.com/peterbourgon/ff/v3"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

func main() {
	mainCtx, cancelMainCtx := context.WithCancel(context.Background())

	fs := flag.NewFlagSet("ghostidp", flag.ContinueOnError)
	port := fs.Int("port", 8080, "port to listen on")
	hydraAdminURL := fs.String("hydra-admin-url", "http://localhost:4445/admin", "hydra admin api url")
	usersFile := fs.String("users-file", "users.yaml", "hard-coded users file")
	debug := fs.Bool("debug", false, "log debug information")
	badge := fs.String("badge", "", "Badge to display")
	version := fs.String("version", "", "Version to display")
	accentColor := fs.String("accent-color", "", "Color to use for accent color")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVars()); err != nil {
		log.Fatal(fmt.Errorf("failed to parse arguments: %v", err))
	}

	// Creating the main logger
	var opts []logger.Opts
	if *debug {
		opts = append(opts, logger.WithLogLevel(logrus.DebugLevel))
	}
	appLogger := logger.New(opts...)

	// Display configuration
	appLogger.Debugf("port: %d", *port)
	appLogger.Debugf("hydra admin URL: %s", *hydraAdminURL)
	appLogger.Debugf("users file: %s", *usersFile)
	appLogger.Debugf("badge: %s", *badge)
	appLogger.Debugf("version: %s", *version)
	appLogger.Debugf("accent color: %s", *accentColor)

	app, cancelAppCtx, err := application.NewApplication(mainCtx, &application.Params{
		Log:           appLogger,
		HydraAdminURL: *hydraAdminURL,
		UsersFile:     *usersFile,
		Display: &application.DisplayParams{
			Version:     *version,
			Badge:       *badge,
			AccentColor: *accentColor,
		},
	})
	if err != nil {
		appLogger.Fatal(fmt.Errorf("failed to create application: %w", err))
	}
	defer func() {
		cancelAppCtx()
		cancelMainCtx()
	}()

	// Create the Wev Server
	router, err := handlers.CreateWebServer(app)
	if err != nil {
		appLogger.Fatal(fmt.Errorf("failed to create Web Server: %v", err))
	}

	// Run the Web Server
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), router); err != nil {
		log.Fatal(err)
	}
}
