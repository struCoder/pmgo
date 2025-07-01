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

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [name]",
		Short: "Stop a process",
		Long:  `Stop a running process by name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			log.Infof("Stopping process %s", name)
			return sendProcessAction(name, "stop")
		},
	}
}

func newRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart [name]",
		Short: "Restart a process",
		Long:  `Restart a process by name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			log.Infof("Restarting process %s", name)
			return sendProcessAction(name, "restart")
		},
	}
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete a process",
		Long:  `Stop and delete a process by name.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			log.Infof("Deleting process %s", name)
			return deleteProcess(name)
		},
	}
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [name]",
		Short: "Show process information",
		Long:  `Display detailed information about a process.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			log.Infof("Showing info for process %s", name)
			return showProcessInfo(name)
		},
	}
}

func newSaveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "Save current process list",
		Long:  `Save the current list of processes to disk.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Saving process list...")
			return saveProcesses()
		},
	}
}

func newKillCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kill",
		Short: "Kill the PMGO daemon",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("Killing PMGO daemon...")
			return nil
		},
	}
}

func newWebCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "Start web interface",
		Long:  `Start the web interface for process management.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			host := viper.GetString("host")
			webPort := viper.GetInt("web-port")
			if webPort == 0 {
				webPort = 8080
			}

			fmt.Printf("Web interface available at: http://%s:%d\n", host, webPort)
			fmt.Println("Press Ctrl+C to exit")

			// Keep the command running
			select {}
		},
	}
}

func newLogsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logs [name]",
		Short: "Show process logs",
		Long:  `Display logs for a specific process.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			log.Infof("Showing logs for process %s", name)
			return showProcessLogs(name)
		},
	}
}

func sendProcessAction(name, action string) error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes/%s/%s", host, port, name, action)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	log.Infof("Process %s %sed successfully", name, action)
	return nil
}

func deleteProcess(name string) error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes/%s", host, port, name)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	log.Infof("Process %s deleted successfully", name)
	return nil
}

func showProcessInfo(name string) error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes/%s", host, port, name)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("process %s not found", name)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	var procInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&procInfo); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	// Display process info
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Property", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for key, value := range procInfo {
		table.Append([]string{key, fmt.Sprintf("%v", value)})
	}

	table.Render()
	return nil
}

func saveProcesses() error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/save", host, port)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	log.Info("Process list saved successfully")
	return nil
}

func showProcessLogs(name string) error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes/%s/logs", host, port, name)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to daemon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("process %s not found", name)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	var logsResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&logsResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	if logs, ok := logsResponse["logs"].(string); ok {
		fmt.Print(logs)
	} else {
		fmt.Println("No logs available")
	}

	return nil
}
