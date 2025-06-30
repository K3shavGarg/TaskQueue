package main

import (
	"Job-Queue/internal/config"
	"Job-Queue/internal/model"
	"Job-Queue/internal/routes"
	"Job-Queue/internal/worker"
	"Job-Queue/metrics"
	"Job-Queue/pkg"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

// The thing is, when you have a black box, you can only ever use it as a black box. But if you have a box that you understand well, you can take it into pieces, repurpose the pieces for completely different things, all without getting lost.

func main() {
	queue := model.NewRedisQueue(config.Client)

	router := mux.NewRouter()
	sm := config.Client.Ping(config.Ctx)
	if sm == nil {
		print("something wrong")
	}
	routes.RegisterJobRoutes(router, queue)
	metrics.Init()
	router.Handle("/metrics", promhttp.Handler())

	// Starting Workers
	worker.StartQueueProcessor(queue, 5, 30)

	pkg.Log.Info(fmt.Sprintf("Listening on PORT: %v", config.Env.Port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.Env.Port), router); err != nil {
		log.Fatal("error in starting the server", err)
	}
}
