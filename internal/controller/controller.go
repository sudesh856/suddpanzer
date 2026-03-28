package controller

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"text/tabwriter"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"

	pb "github.com/sudesh856/suddpanzer/proto"
)

type AgentConfig struct {
	ID      string `yaml:"id"`
	Address string `yaml:"address"`
	Region  string `yaml:"region"`
}

type AgentsFile struct {
	Agents []AgentConfig `yaml:"agents"`
}

type AgentResult struct {
	AgentID       string
	Region        string
	TotalRequests int64
	ErrorCount    int64
	AvgRPS        float64
	P50Ms         int64
	P95Ms         int64
	P99Ms         int64
	ErrorRatePct  float64
}

func LoadAgentsFile(path string) (*AgentsFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Read agents file %q: %w", path, err)
	}
	var af AgentsFile
	if err := yaml.Unmarshal(data, &af); err != nil {
		return nil, fmt.Errorf("Parse agents yaml: %w", err)
	}
	if len(af.Agents) == 0 {
		return nil, fmt.Errorf("Agents.yaml has no agents")
	}
	return &af, nil
}

func Run(ctx context.Context, af *AgentsFile, scenarioYAML string) ([]AgentResult, error) {
	type entry struct {
		cfg    AgentConfig
		conn   *grpc.ClientConn
		stream pb.SuddpanzerAgent_RunScenarioClient
	}

	entries := make([]entry, 0, len(af.Agents))

	for _, cfg := range af.Agents {
		fmt.Printf("[Controller] dialing %s at %s\n", cfg.ID, cfg.Address)
		conn, err := grpc.Dial(cfg.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("dial agent %s: %w", cfg.ID, err)
		}
		entries = append(entries, entry{cfg: cfg, conn: conn})
	}
	defer func() {
		for _, e := range entries {
			e.conn.Close()
		}
	}()

	for i := range entries {
		client := pb.NewSuddpanzerAgentClient(entries[i].conn)
		stream, err := client.RunScenario(ctx, &pb.WorkSpec{
			ScenarioYaml: scenarioYAML,
			AgentId:      entries[i].cfg.ID,
			Region:       entries[i].cfg.Region,
		})
		if err != nil {
			return nil, fmt.Errorf("Start scenario on %s: %w", entries[i].cfg.ID, err)
		}
		entries[i].stream = stream
		fmt.Printf("[Controller] agent %s started\n", entries[i].cfg.ID)
	}
	fmt.Println()

	type tagged struct {
		agentID string
		snap    *pb.MetricsSnapshot
	}

	snapCh := make(chan tagged, 256)
	var wg sync.WaitGroup

	for _, e := range entries {
		e := e
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				snap, err := e.stream.Recv()
				if err == io.EOF {
					fmt.Printf("\n[controller] agent %s finished\n", e.cfg.ID)
					return
				}
				if err != nil {
					fmt.Printf("\n[controller] agent %s error: %v\n", e.cfg.ID, err)
					return
				}
				snapCh <- tagged{agentID: e.cfg.ID, snap: snap}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(snapCh)
	}()

	latest := make(map[string]*pb.MetricsSnapshot)
	var mu sync.Mutex
	start := time.Now()

	printTicker := time.NewTicker(1 * time.Second)
	defer printTicker.Stop()

	doneCh := make(chan struct{})
	go func() {
		for snap := range snapCh {
			mu.Lock()
			latest[snap.agentID] = snap.snap
			mu.Unlock()
		}
		close(doneCh)
	}()

	for {
		select {
		case <-printTicker.C:
			mu.Lock()
			printLive(latest, time.Since(start))
			mu.Unlock()
		case <-doneCh:
			mu.Lock()
			results := buildResults(latest)
			mu.Unlock()
			fmt.Println()
			printFinalTable(results)
			return results, nil
		}
	}
}

func printLive(latest map[string]*pb.MetricsSnapshot, elapsed time.Duration) {
	if len(latest) == 0 {
		return
	}
	var totalRPS float64
	var totalReqs, totalErr int64
	line := fmt.Sprintf("\r[%4.0fs]", elapsed.Seconds())
	for id, s := range latest {
		line += fmt.Sprintf(" | %s RPS:%.1f p99:%dms err:%.1f%%", id, s.Rps, s.P99Ms, s.ErrorRatePct)
		totalRPS += s.Rps
		totalReqs += s.TotalRequests
		totalErr += s.ErrorCount
	}
	var errRate float64
	if totalReqs > 0 {
		errRate = float64(totalErr) / float64(totalReqs) * 100
	}
	line += fmt.Sprintf(" | TOTAL RPS:%.1f err:%.1f%%  ", totalRPS, errRate)
	fmt.Print(line)
}

func buildResults(latest map[string]*pb.MetricsSnapshot) []AgentResult {
	out := make([]AgentResult, 0, len(latest))
	for _, s := range latest {
		out = append(out, AgentResult{
			AgentID:       s.AgentId,
			Region:        s.Region,
			TotalRequests: s.TotalRequests,
			ErrorCount:    s.ErrorCount,
			AvgRPS:        s.Rps,
			P50Ms:         s.P50Ms,
			P95Ms:         s.P95Ms,
			P99Ms:         s.P99Ms,
			ErrorRatePct:  s.ErrorRatePct,
		})
	}
	return out
}

func printFinalTable(results []AgentResult) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("  DISTRIBUTED RUN - PER-AGENT BREAKDOWN")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "  Agent\tRegion\tRequests\tRPS\tp50ms\tp95ms\tp99ms\tErrors\tErrRate%")
	fmt.Fprintln(w, "  ─────\t──────\t────────\t───\t─────\t─────\t─────\t──────\t────────")

	var totalReqs, totalErr int64
	var totalRPS float64

	for _, r := range results {
		fmt.Fprintf(w, "  %s\t%s\t%d\t%.1f\t%d\t%d\t%d\t%d\t%.2f\n",
			r.AgentID, r.Region, r.TotalRequests, r.AvgRPS,
			r.P50Ms, r.P95Ms, r.P99Ms, r.ErrorCount, r.ErrorRatePct)
		totalReqs += r.TotalRequests
		totalErr += r.ErrorCount
		totalRPS += r.AvgRPS
	}

	var totalErrRate float64
	if totalReqs > 0 {
		totalErrRate = float64(totalErr) / float64(totalReqs) * 100
	}

	fmt.Fprintln(w, "  ─────\t──────\t────────\t───\t─────\t─────\t─────\t──────\t────────")
	fmt.Fprintf(w, "  TOTAL\t\t%d\t%.1f\t\t\t\t%d\t%.2f\n",
		totalReqs, totalRPS, totalErr, totalErrRate)
	w.Flush()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
