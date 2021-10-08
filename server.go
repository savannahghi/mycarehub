package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/savannahghi/serverutils"
	log "github.com/sirupsen/logrus"
	"go.opencensus.io/stats/view"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation"
)

const waitSeconds = 30

func main() {
	ctx := context.Background()
	err := serverutils.Sentry()
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	// Firstly, we'll register Server views.
	// A service can declare it's own additional views
	if err := view.Register(serverutils.DefaultServiceViews...); err != nil {
		serverutils.LogStartupError(ctx, err)
	}

	deferFunc, err := serverutils.EnableStatsAndTraceExporters(ctx, serverutils.MetricsCollectorService("onboarding"))
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	defer deferFunc()

	// initialize the tracing provider in prod and testing env only
	// env := serverutils.GetRunningEnvironment()
	// if env == serverutils.ProdEnv || env == serverutils.TestingEnv {
	// 	tp, err := serverutils.InitOtelSDK(ctx, "onboarding")
	// 	if err != nil {
	// 		serverutils.LogStartupError(ctx, err)
	// 	}
	// 	defer tp.Shutdown(ctx)
	// }

	port, err := strconv.Atoi(serverutils.MustGetEnvVar(serverutils.PortEnvVarName))
	if err != nil {
		serverutils.LogStartupError(ctx, err)
	}
	srv := presentation.PrepareServer(ctx, port, presentation.AllowedOrigins)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			serverutils.LogStartupError(ctx, err)
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
