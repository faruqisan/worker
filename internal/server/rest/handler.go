package rest

import (
	"log"
	"net/http"
	"time"

	"github.com/faruqisan/worker/internal/services/jobs"
)

// HandlePushHTTPJob ..
func (s Server) HandlePushHTTPJob(w http.ResponseWriter, r *http.Request) {
	job := jobs.HTTPRequestJob{
		URL:      "http://faruqisan.com",
		Method:   http.MethodGet,
		Schedule: time.Now().Add(time.Minute),
	}

	err := s.jobSVC.PushHTTPJob(job)
	if err != nil {
		log.Println(err)
		w.Write([]byte("err"))
	}

	w.Write([]byte("sucess"))
}
