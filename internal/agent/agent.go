package agent

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
 
	"google.golang.org/grpc"
 
	pb "github.com/sudesh856/suddpanzer/proto"
	"github.com/sudesh856/suddpanzer/internal/metrics"
	"github.com/sudesh856/suddpanzer/internal/pool"
	"github.com/sudesh856/suddpanzer/internal/ramp"
	"github.com/sudesh856/suddpanzer/internal/scenario"
	"github.com/sudesh856/suddpanzer/internal/worker"
)

type Server struct {
	pb.UnimplementedSuddpanzerAgentServer
	id string
	region string
	gs *grpc.Server
}
func New(id, region string) *Server {
	return &Server{id: id, region: region}
}

func(s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("Agent listen %s: %w", addr, err)
	}

	s.gs = grpc.NewServer()
	pb.RegisterSuddpanzerAgentServer(s.gs, s)
	fmt.Printf("[Agent:%s] listening on %s (region: %s)\n", s.id, addr, s.region)
	return s.gs.Serve(lis)
}

func(s *Server) Stop() {
	if s.gs != nil {
		s.gs.GracefulStop()
	}
}
func (s *Server) RunScenario(spec *pb.WorkSpec, stream pb.SuddpanzerAgent_RunScenarioServer) error {
	fmt.Printf("[Agent:%s] received WorkSpec — starting\n", s.id)
 
	tmp, err := os.CreateTemp("", "suddpanzer-*.yaml")
	if err != nil {
		return fmt.Errorf("Temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.WriteString(spec.ScenarioYaml); err != nil {
		return fmt.Errorf("Write yaml: %w", err)
	}
	tmp.Close()
 
	sc, err := scenario.LoadScenario(tmp.Name())
	if err != nil {
		return fmt.Errorf("Parse scenario: %w", err)
	}

	rampStages := make([]ramp.Stage, len(sc.Stages))
	totalDur := time.Duration(0)
	maxVUs := 0
	for i, st := range sc.Stages {
		rampStages[i] = ramp.Stage{Duration: st.ParsedDuration, TargetVUs: st.TargetVUs}
		totalDur += st.ParsedDuration
		if st.TargetVUs > maxVUs {
			maxVUs = st.TargetVUs
		}
	}
 
	ctx, cancel := context.WithTimeout(stream.Context(), totalDur+5*time.Second)
	defer cancel()
 
	ctrl := ramp.New(rampStages)
	go ctrl.Run(ctx)
 
	p := pool.New(1000)
	p.Start(ctx, maxVUs)
	agg := metrics.New()
	agg.Start(p.Results())
	start := time.Now()
 
	chainStore := scenario.NewChainStore()
	go func() {
		for result := range p.Results() {
			if result.EndpointName != "" && len(result.Body) > 0 {
				for _, ep := range sc.Endpoints {
					if ep.Name == result.EndpointName {
						chainStore.Store(result.EndpointName, result.Body, ep.Extract)
						break
					}
				}
			}
		}
	}()
 
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case ctrl.Semaphore <- struct{}{}:
				ep := scenario.PickEndpoint(sc.Endpoints)
				vars := chainStore.ToVars()
				p.Submit(worker.Job{
				Name:           ep.Name,
				URL:            scenario.ReplaceVars(ep.URL, vars),
				Method:         ep.Method,
				Body:           scenario.ReplaceVars(ep.Body, vars),
				ExpectedStatus: ep.ExpectedStatus,
				Headers:        ep.Headers,
				BasicAuth:      ep.BasicAuth,
			})
			select {
			case <-ctrl.Semaphore:
			case <-ctx.Done():
				return
			}
			}
		}
	}()
 

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
 
	snap := func() *pb.MetricsSnapshot {
		elapsed := time.Since(start)
		return &pb.MetricsSnapshot{
			AgentId:       s.id,
			Region:        s.region,
			Timestamp:     time.Now().Unix(),
			Rps:           agg.RPS(elapsed),
			P50Ms:         agg.P50(),
			P95Ms:         agg.P95(),
			P99Ms:         agg.P99(),
			ErrorRatePct:  agg.ErrorRate(),
			TotalRequests: agg.TotalRequests(),
			ErrorCount:    agg.ErrorCount(),
			VusActive:     int64(maxVUs),
		}
	}
 
	for {
		select {
		case <-ctx.Done():
			_ = stream.Send(snap())
			wg.Wait()
			fmt.Printf("[Agent:%s] done\n", s.id)
			return nil
		case <-ticker.C:
			s := snap()
			fmt.Printf("[Agent:%s] RPS:%.1f p99:%dms err:%.2f%% total:%d\n",
				s.AgentId, s.Rps, s.P99Ms, s.ErrorRatePct, s.TotalRequests)
			if err := stream.Send(s); err != nil {
				return err
			}
		}
	}
}
 