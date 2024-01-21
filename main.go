package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go.dsage.org/gopkgs/internal"
)

func main() {
	// create and activate a new text handler w/ slog
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	// define the flags and parse for them
	help := flag.Bool("help", false, "show this usage information")
	port := flag.Int("port", 8080, "the port that the server should listen on")
	flag.Parse()

	// check if the user requested the command line parameters
	if *help {
		flag.Usage()
		return
	}

	// print the project name and version
	slog.Info("starting gopkgs", "version", BuildID, "copyright", "(c) 2024 Daniel Sage (dsage.org); see LICENSE")

	// open a signal channel and listen for SIGTERM interrupts
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt, syscall.SIGTERM)

	// start the core listeners in goroutines
	go internal.Listen(fmt.Sprintf(":%d", *port))

	// block the main thread, waiting for an interrupt
	<-sc
	slog.Info("interrupt received, stopping the server")

	// stop the core listeners
	internal.Shutdown()
}
