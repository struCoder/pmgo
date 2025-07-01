package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type StartRequest struct {
	Command string   `json:"command"`
	Name    string   `json:"name"`
	Args    []string `json:"args,omitempty"`
	Binary  bool     `json:"binary,omitempty"`
}

func newStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start [source] [name]",
		Short: "Start a new process",
		Long:  `Start and daemonize a new Go application process.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sourcePath := args[0]
			name := args[1]

			// Get flags
			processArgs, _ := cmd.Flags().GetStringSlice("args")
			binary, _ := cmd.Flags().GetBool("binary")

			log.Infof("Starting process %s from %s", name, sourcePath)

			// Create request
			request := StartRequest{
				Command: sourcePath,
				Name:    name,
				Args:    processArgs,
				Binary:  binary,
			}

			// Send HTTP request to daemon
			return sendStartRequest(request)
		},
	}

	cmd.Flags().StringSlice("args", []string{}, "Arguments to pass to the process")
	cmd.Flags().Bool("binary", false, "Start from compiled binary")

	return cmd
}

func sendStartRequest(request StartRequest) error {
	host := viper.GetString("host")
	port := viper.GetInt("port")

	url := fmt.Sprintf("http://%s:%d/api/v1/processes", host, port)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("daemon returned status: %d", resp.StatusCode)
	}

	log.Infof("Process started successfully")
	return nil
}
