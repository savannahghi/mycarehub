package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation"
)

const waitSeconds = 30

func main() {
	ctx := context.Background()
	err := base.Sentry()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	port, err := strconv.Atoi(base.MustGetEnvVar(base.PortEnvVarName))
	if err != nil {
		base.LogStartupError(ctx, err)
	}
	srv := presentation.PrepareServer(ctx, port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			base.LogStartupError(ctx, err)
		}
	}()

	// Block until we receive a sigint (CTRL+C) signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*waitSeconds)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until timeout
	err = srv.Shutdown(ctx)
	log.Printf("graceful shutdown started; the timeout is %d secs", waitSeconds)
	if err != nil {
		log.Printf("error during clean shutdown: %s", err)
		os.Exit(-1)
	}
	os.Exit(0)
}
