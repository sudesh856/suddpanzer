package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sudesh856/LoadForge/internal/metrics"
	"github.com/sudesh856/LoadForge/internal/pool"
	"github.com/sudesh856/LoadForge/internal/worker"
	"github.com/spf13/cobra"
	"golang.org/x/time/rate"
)

var url      string
var vus      int
var duration string
var rps      int

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a load test",
	Run: func(cmd *cobra.Command, args []string) {
		dur, err := time.ParseDuration(duration)
		if err != nil {
			fmt.Println("invalid duration:", err)
			return
		}

		// ctx cancelled on timeout OR ctrl+C
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		ctx, cancel = context.WithTimeout(ctx, dur)
		defer cancel()

		// rate limiter
		var limiter *rate.Limiter
		if rps > 0 {
			limiter = rate.NewLimiter(rate.Limit(rps), rps)
		}

		p := pool.New(1000)
		p.Start(ctx, vus)

		agg := metrics.New()
		agg.Start(p.Results())

		start := time.Now()

		// WaitGroup to track job submitter goroutine
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					if limiter != nil {
						limiter.Wait(ctx)
					}
					p.Submit(worker.Job{URL: url})
				}
			}
		}()

		// live terminal output every second
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					elapsed := time.Since(start)
					fmt.Printf("\rRequests: %d | RPS: %.0f | p99: %dms | Errors: %d",
						agg.TotalRequests(),
						agg.RPS(elapsed),
						agg.P99(),
						agg.ErrorCount(),
					)
				}
			}
		}()

		// wait for context to expire or ctrl+C
		<-ctx.Done()

		// wait for submitter goroutine to stop
		wg.Wait()

		// give aggregator a moment to process remaining results
		time.Sleep(100 * time.Millisecond)

		elapsed := time.Since(start)

		// final summary
		fmt.Printf("\n\n---- SUDD LOAD TEST SUMMARY ----\n")
		fmt.Printf("URL            : %s\n", url)
		fmt.Printf("VUs            : %d\n", vus)
		fmt.Printf("Duration       : %s\n", elapsed.Round(time.Second))
		fmt.Printf("-----------------------------------\n")
		fmt.Printf("Total Requests : %d\n", agg.TotalRequests())
		fmt.Printf("Avg RPS        : %.2f\n", agg.RPS(elapsed))
		fmt.Printf("-----------------------------------\n")
		fmt.Printf("p50            : %dms\n", agg.P50())
		fmt.Printf("p75            : %dms\n", agg.P75())
		fmt.Printf("p90            : %dms\n", agg.P90())
		fmt.Printf("p95            : %dms\n", agg.P95())
		fmt.Printf("p99            : %dms\n", agg.P99())
		fmt.Printf("p999           : %dms\n", agg.P999())
		fmt.Printf("Max            : %dms\n", agg.Max())
		fmt.Printf("-----------------------------------\n")
		fmt.Printf("Errors         : %d\n", agg.ErrorCount())
		fmt.Printf("Error Rate     : %.2f%%\n", agg.ErrorRate())
		fmt.Printf("===================================\n")
	},
}

func init() {
	runCmd.Flags().StringVar(&url,      "url",      "",    "Target URL")
	runCmd.Flags().IntVar(&vus,         "vus",      10,    "Virtual users")
	runCmd.Flags().StringVar(&duration, "duration", "30s", "Test duration")
	runCmd.Flags().IntVar(&rps,         "rps",      0,     "Max requests per second (0 = unlimited)")
	rootCmd.AddCommand(runCmd)
}