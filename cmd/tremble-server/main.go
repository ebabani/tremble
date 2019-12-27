package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ebabani/tremble/twitch"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var clientID string

func startServer(port string, r chi.Router) chan bool {
	exit := make(chan bool)

	go func() {
		zap.S().Infof("Starting server on port %s\n", port)
		log.Println(http.ListenAndServe(port, r))
		zap.S().Infof("Shutting down server on port %s\n", port)
		exit <- true
	}()
	return exit
}

func setupLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)

	zap.RedirectStdLog(logger)
	return logger, nil
}

func main() {
	logger, err := setupLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	log.Println("Starting tremble")
	clientID = os.Getenv("TWITCH_CLIENT_ID")

	appRouter := chi.NewRouter()

	metricsRouter := chi.NewRouter()

	appRouter.Get("/videos", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strings.Join(twitch.GetVideos(""), ",")))
	})

	select {
	case <-startServer("127.0.0.1:8080", appRouter):
		log.Println("App Server Shut Down")
	case <-startServer("127.0.0.1:8081", metricsRouter):
		fmt.Println("Metrics Server Shut Down")
	}
}
