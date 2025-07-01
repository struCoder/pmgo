package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ProcessInfo struct {
	Name        string   `json:"name"`
	PID         int      `json:"pid"`
	Status      string   `json:"status"`
	Uptime      string   `json:"uptime,omitempty"`
	Restarts    int      `json:"restarts,omitempty"`
	CPUPercent  float64  `json:"cpu_percent,omitempty"`
	MemoryBytes int64    `json:"memory_bytes,omitempty"`
	Command     string   `json:"command,omitempty"`
	Args        []string `json:"args,omitempty"`
}

type ProcessListResponse struct {
	Processes []ProcessInfo `json:"processes"`
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "status"},
		Short:   "List all managed processes",
		Long:    `Display status information for all managed processes.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Listing all processes...")
			return listProcesses()
		},
	}

	return cmd
}

func listProcesses() error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes", host, port)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	var response ProcessListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Display table
	if len(response.Processes) == 0 {
		fmt.Println("No processes are currently managed by PMGO.")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "PID", "Status", "Uptime", "Restarts", "CPU", "Memory"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetBorder(true)

	for _, proc := range response.Processes {
		// Format CPU percentage - always show percentage, even for 0
		cpuStr := fmt.Sprintf("%.1f%%", proc.CPUPercent)

		// Format memory
		memoryStr := "-"
		if proc.MemoryBytes > 0 {
			memoryStr = formatMemory(proc.MemoryBytes)
		}

		// Format uptime
		uptimeStr := proc.Uptime
		if uptimeStr == "" {
			uptimeStr = "-"
		}

		row := []string{
			proc.Name,
			fmt.Sprintf("%d", proc.PID),
			proc.Status,
			uptimeStr,
			fmt.Sprintf("%d", proc.Restarts),
			cpuStr,
			memoryStr,
		}
		table.Append(row)
	}

	table.Render()
	return nil
}

// formatMemory formats memory bytes in human-readable format
func formatMemory(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	if bytes >= GB {
		return fmt.Sprintf("%.1fGB", float64(bytes)/GB)
	} else if bytes >= MB {
		return fmt.Sprintf("%.1fMB", float64(bytes)/MB)
	} else if bytes >= KB {
		return fmt.Sprintf("%.1fKB", float64(bytes)/KB)
	}
	return fmt.Sprintf("%dB", bytes)
}
