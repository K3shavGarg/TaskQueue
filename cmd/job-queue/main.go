package main

import (
	"Job-Queue/internal/config"
	"Job-Queue/internal/routes"
	"Job-Queue/internal/worker"
	"Job-Queue/pkg"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// The thing is, when you have a black box, you can only ever use it as a black box. But if you have a box that you understand well, you can take it into pieces, repurpose the pieces for completely different things, all without getting lost.

func main() {
	router := mux.NewRouter()

	routes.RegisterJobRoutes(router)
	// Starting Workers
	worker.StartWorkerPool(routes.JobService.Queue, 6)

	pkg.Log.Info(fmt.Sprintf("Listening on PORT: %v", config.Env.Port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", config.Env.Port), router); err != nil {
		log.Fatal("error in starting the server", err)
	}
}
