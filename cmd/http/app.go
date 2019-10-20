package main

import (
	"log"

	jobsrsc "github.com/faruqisan/worker/internal/resources/jobs"
	"github.com/faruqisan/worker/internal/server/rest"
	"github.com/faruqisan/worker/internal/services/jobs"
	"github.com/faruqisan/worker/pkg/cache"
)

func main() {

	cache := cache.New()
	jobsRSC := jobsrsc.New(cache)
	jobsSVC := jobs.New(jobsRSC)
	rest := rest.New(jobsSVC)

	log.Println("server is running")
	log.Fatal(rest.Run())
}
