package ramp
import (
	"context"
	"time"
)

type Controller struct {
	Semaphore chan struct {}
	stages []Stage
}

type Stage struct {
	Duration time.Duration
	TargetVUs int
}

func New(stages []Stage) *Controller {
	return &Controller{
		Semaphore: make(chan struct{}, maxVUs(stages)),
		stages: stages,
	}
}

func(c *Controller) Run(ctx context.Context) {
	currentVUs := 0

	for _, stage := range c.stages {
		startVUs := currentVUs
		targetVUs := stage.TargetVUs
		duration := stage.Duration

		if duration == 0 {
			c.setVUs(targetVUs)
			currentVUs = targetVUs
			continue
		}


		ticker := time.NewTicker(100 * time.Millisecond)
		stageStart := time.Now()

		stageLoop:
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				elapsed := time.Since(stageStart)
				if elapsed >= duration {
					c.setVUs(targetVUs)
					currentVUs = targetVUs
					ticker.Stop()
					break stageLoop
				}
				//linear interpolation
				progress := float64(elapsed) / float64(duration)
				vus := startVUs + int(float64(targetVUs - startVUs)*progress)
				c.setVUs(vus)
			}
		}
	}
}

func(c *Controller) setVUs(n int) {
	current := len(c.Semaphore)

	switch {
	case n > current:
		for i := 0; i <n-current; i++ {
			select {
			case c.Semaphore <- struct{}{}:
			default:

			}
		}

	case n < current:
		for i :=0; i < current-n; i++ {
			select {
			case <-c.Semaphore:
			default:

			}
		}
	}
}


func maxVUs(stages []Stage) int {
	max := 0
	for _, s := range stages {
		if s.TargetVUs > max {
			max = s.TargetVUs
		}
	}

	if max == 0 {
		max = 1
	}
	return  max
}