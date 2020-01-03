package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/ebabani/tremble/logger"
	"github.com/ebabani/tremble/tracer"
	"github.com/ebabani/tremble/twitch"
	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	"github.com/opentracing-contrib/go-zap/log"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var clientID string

func startServer(port string, r chi.Router) chan bool {
	exit := make(chan bool)

	go func() {
		zap.S().Infof("Starting server on port %s\n", port)
		err := http.ListenAndServe(port, r)
		if err != nil {
			log.Info(err.Error())
		}
		zap.S().Infof("Shutting down server on port %s\n", port)
		exit <- true
	}()
	return exit
}

func main() {
	logger, err := logger.SetupLogger()
	if err != nil {
		log.Fatal(err.Error())
	}
	defer logger.Sync()

	closer, err := tracer.SetupTracer("tremble-server")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer closer.Close()

	log.Info("Starting tremble")
	clientID = os.Getenv("TWITCH_CLIENT_ID")

	appRouter := chi.NewRouter()
	metricsRouter := chi.NewRouter()
	twitchClinet := twitch.TwitchClient{}

	appRouter.Use(func(next http.Handler) http.Handler {
		return nethttp.Middleware(opentracing.GlobalTracer(), next)
	})

	metricsRouter.Use(func(next http.Handler) http.Handler {
		return nethttp.Middleware(opentracing.GlobalTracer(), next)
	})

	appRouter.Get("/videos", func(w http.ResponseWriter, r *http.Request) {
		// span, ctx := opentracing.StartSpanFromContext(r.Context(), "/videos")
		// defer span.Finish()
		w.Write([]byte(strings.Join(twitchClinet.GetVideos(r.Context(), ""), ",")))
	})

	metricsRouter.Get("/metrics", promhttp.Handler().ServeHTTP)
	select {
	case <-startServer("127.0.0.1:8080", appRouter):
		log.Info("App Server Shut Down")
	case <-startServer("127.0.0.1:8081", metricsRouter):
		log.Info("Metrics Server Shut Down")
	}

	lightstep.Flush(context.Background(), opentracing.GlobalTracer())
}
