package jobs

import (
	"encoding/json"

	jobssvc "github.com/faruqisan/worker/internal/services/jobs"
	"github.com/faruqisan/worker/pkg/cache"
	"github.com/go-redis/redis"
)

type (
	// Resource struct ..
	Resource struct {
		cache cache.Engine
	}
)

const (
	httpJobListKey = "queue:http-job"
)

// New ..
func New(cache cache.Engine) Resource {
	return Resource{
		cache: cache,
	}
}

// PushHTTPJob ..
func (r Resource) PushHTTPJob(job jobssvc.HTTPRequestJob) error {

	raw, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = r.cache.ZAdd(httpJobListKey, redis.Z{
		Score:  float64(job.Schedule.Unix()),
		Member: string(raw),
	}).Result() // , ).Err()

	return err
}

// FetchHTTPJobs ..
func (r Resource) FetchHTTPJobs() ([]jobssvc.HTTPRequestJob, error) {
	var (
		jobs    []jobssvc.HTTPRequestJob
		rawJobs []string
		err     error
	)

	rawJobs, err = r.cache.ZRangeByScore(httpJobListKey, redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()
	if err != nil {
		return jobs, err
	}

	for _, rawJob := range rawJobs {
		var (
			j jobssvc.HTTPRequestJob
		)

		err = json.Unmarshal([]byte(rawJob), &j)
		if err != nil {
			return jobs, err
		}

		jobs = append(jobs, j)

	}

	return jobs, err

}

// RemoveHTTPJob ..
func (r Resource) RemoveHTTPJob(job jobssvc.HTTPRequestJob) error {
	raw, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.cache.ZRem(httpJobListKey, raw).Err()
}
