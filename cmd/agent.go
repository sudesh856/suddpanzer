package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sudesh856/suddpanzer/internal/agent"
)

var agentID     string
var agentAddr   string
var agentRegion string

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start a distributed load test agent",
	Long: `Start a suddpanzer agent. The agent waits for a controller to send
a scenario, runs it locally, and streams metrics back every second.

Examples:
  suddpanzer agent --id agent-1 --addr :7071
  suddpanzer agent --id agent-2 --addr :7072 --region eu-west
  suddpanzer agent --id agent-vps --addr :7071 --region us-east`,

	Run: func(cmd *cobra.Command, args []string) {
		if agentID == "" {
			fmt.Println("error: --id is required")
			os.Exit(1)
		}

		srv := agent.New(agentID, agentRegion)

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-sigCh
			fmt.Printf("\n[agent:%s] shutting down\n", agentID)
			srv.Stop()
		}()

		if err := srv.Start(agentAddr); err != nil {
			fmt.Printf("[agent:%s] error: %v\n", agentID, err)
			os.Exit(1)
		}
	},
}

func init() {
	agentCmd.Flags().StringVar(&agentID,     "id",     "",      "Agent ID (required)")
	agentCmd.Flags().StringVar(&agentAddr,   "addr",   ":7071", "Address to listen on")
	agentCmd.Flags().StringVar(&agentRegion, "region", "",      "Region label e.g. us-east")
	rootCmd.AddCommand(agentCmd)
}