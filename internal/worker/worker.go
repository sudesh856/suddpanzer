package worker

import (
	"context"
	"net/http"
	"time"
	"io"
	"strings"
)

//since Job is what we send to a worker which is  the url to hit;

type Job struct {
	URL string
	Method string
	Body string
}


//and since Result is what the worker sends back after firing a request

type Result struct {
	Latency time.Duration
	StatusCode int
	Err error
	Bytes int64
}

//listening by RunWorker for jobs and firing HTTP requests

func RunWorker(ctx context.Context, jobs <- chan Job, results chan <- Result) {

	client := &http.Client{Timeout: 10 * time.Second}

	for {
		select {
		case <-ctx.Done():
			return //the context is now cancelled, stop worker.
		case job,ok := <-jobs:
			if !ok {
				return //job channel is now closed, stop worker.
			}

			start := time.Now()

			method := job.Method
			if method == "" {
				method = "GET"
			}
			var bodyReader io.Reader
			if job.Body != "" {
				bodyReader = strings.NewReader(job.Body)
			}
			req, err := http.NewRequestWithContext(ctx, method, job.URL, bodyReader)
			if err != nil {
				results <- Result{Latency: time.Since(start), Err: err}
				continue
			}
			resp, err := client.Do(req)

			latency := time.Since(start)

			if err != nil {
				results <- Result {Latency: latency, Err: err}
				continue }

					results <- Result{
						Latency: latency,
						StatusCode: resp.StatusCode,
						Bytes: resp.ContentLength,
					}
					resp.Body.Close()
				}

				
			}
		}
	