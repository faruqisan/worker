package jobs

import (
	"log"
	"net/http"
	"time"
)

type (
	// HTTPRequestJob struct define job for http request
	HTTPRequestJob struct {
		URL    string            `json:"url"`
		Method string            `json:"method"`
		Header map[string]string `json:"headers"`
		Body   interface{}       `json:"body"`

		Schedule   time.Time `json:"schedule"`
		AllowRetry bool      `json:"allow_retry"`
	}

	resource interface {
		PushHTTPJob(job HTTPRequestJob) error
		FetchHTTPJobs() ([]HTTPRequestJob, error)
		RemoveHTTPJob(job HTTPRequestJob) error
	}

	// Service ..
	Service struct {
		resource resource
	}
)

// New ..
func New(resource resource) Service {
	s := Service{
		resource: resource,
	}
	go func() {
		s.ConsumeJob()
	}()
	return s
}

// PushHTTPJob ..
func (s Service) PushHTTPJob(job HTTPRequestJob) error {
	return s.resource.PushHTTPJob(job)
}

func (s Service) performHTTPJob(j HTTPRequestJob) error {
	log.Println("performing job : ", j)
	resp, err := http.Get(j.URL)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(resp)
	return nil
}

// ConsumeJob ..
func (s Service) ConsumeJob() {

	tick := time.NewTicker(time.Second)

	for {
		select {
		case <-tick.C:
			jobs, err := s.resource.FetchHTTPJobs()
			if err != nil {
				log.Println(err)
				continue
			}

			for _, job := range jobs {

				tN := time.Now()
				if job.Schedule.After(tN) {
					// todo: check idempotency
					// so we don't remove the job before it completed
					// inside go routine
					tDiff := job.Schedule.Sub(tN)
					go func(j HTTPRequestJob) {
						log.Println("performing job inside go routine")
						time.Sleep(tDiff)
						s.performHTTPJob(j)
					}(job)
					s.resource.RemoveHTTPJob(job)
					continue
				}

				log.Println("performing job outside go routine")
				s.performHTTPJob(job)
				s.resource.RemoveHTTPJob(job)
			}
		}
	}
}
