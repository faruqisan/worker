package rest

import (
	"net/http"

	"github.com/faruqisan/worker/internal/services/jobs"
)

type (
	jobService interface {
		PushHTTPJob(job jobs.HTTPRequestJob) error
	}

	// Server struct
	Server struct {
		jobSVC jobService
	}
)

// New ..
func New(jobSVC jobService) Server {
	return Server{
		jobSVC: jobSVC,
	}
}

// Run ..
func (s Server) Run() error {
	return http.ListenAndServe(":8080", s.handle())
}
