package main

// This file contains helpers that initialize app insight, developer tooling and database set up that might be run on any given app
import (
	"fmt"
	"log"
	"strconv"

	"github.com/caring/go-packages/pkg/grpc_middleware"
	"github.com/caring/go-packages/pkg/logging"
	"github.com/caring/go-packages/pkg/tracing"
	"github.com/getsentry/sentry-go"

	"google.golang.org/grpc"
)

// establish logging from env config
func initLogger() *logging.Logger {
	log.Print("Initializing logger")
	l, err := logging.NewLogger(&logging.Config{})
	if err != nil {
		sentry.CaptureException(err)
		log.Fatal("Error initializing logger:", err.Error())
	}
	log.Print("Done")
	return l
}

// configure sentry from env
func initSentry(logger *logging.Logger) {
	logger.Debug("Initializing Sentry")
	val := envMust("SENTRY_DISABLE")
	disabled, err := strconv.ParseBool(val)
	if err != nil {
		logger.Fatal("Error getting SENTRY_DISABLE variable")
	}
	if disabled {
		logger.Debug("Skipping")
		return
	}

	sentryDsn := envMust("SENTRY_DSN")
	env := envMust("SENTRY_ENV")
	err = sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Environment: env,
	})
	if err != nil {
		logger.Fatal("sentry.Init:" + err.Error())
	}
	logger.Debug("Done")

}

// configure tracing form env
func initTracing(logger *logging.Logger) *tracing.Tracer {
	logger.Debug("Initializing Tracing")
	tracer, err := tracing.NewTracer(&tracing.Config{
		Logger: logger,
	})
	if err != nil {
		sentry.CaptureException(err)
		logger.Fatal("Failed to establish tracing:" + err.Error())
	}
	logger.Debug("Done")
	return tracer
}

// create protocol server with chained interceptors
func createGRPCServer(logger *logging.Logger, tracer *tracing.Tracer) *grpc.Server {
	return grpc.NewServer(
		grpc_middleware.NewGRPCChainedUnaryInterceptor(grpc_middleware.UnaryOptions{
			Logger: logger,
			Tracer: tracer,
		}),
		grpc_middleware.NewGRPCChainedStreamInterceptor(grpc_middleware.StreamOptions{
			Logger: logger,
			Tracer: tracer,
		}),
	)
}


