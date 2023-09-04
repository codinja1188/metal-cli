package eventsprojtest

import (
	"io"
	"os"
	"strings"
	"testing"

	root "github.com/equinix/metal-cli/internal/cli"
	"github.com/equinix/metal-cli/internal/events"
	outputPkg "github.com/equinix/metal-cli/internal/outputs"
	"github.com/equinix/metal-cli/test/helper"
	"github.com/spf13/cobra"
)

func TestCli_Events_Get(t *testing.T) {
	var projectId, deviceId string
	var err error
	var status bool
	subCommand := "event"
	consumerToken := ""
	apiURL := ""
	Version := "metal"
	rootClient := root.NewClient(consumerToken, apiURL, Version)
	type fields struct {
		MainCmd  *cobra.Command
		Outputer outputPkg.Outputer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *cobra.Command
		cmdFunc func(*testing.T, *cobra.Command)
	}{
		{
			name: "get_events_by_dev_id",
			fields: fields{
				MainCmd:  events.NewClient(rootClient, outputPkg.Outputer(&outputPkg.Standard{})).NewCommand(),
				Outputer: outputPkg.Outputer(&outputPkg.Standard{}),
			},
			want: &cobra.Command{},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				root := c.Root()
				projectId, err = helper.CreateTestProject("metal-cli-events-pro")
				if err != nil {
					t.Error(err)
				}
				deviceId, err = helper.CreateTestDevice(projectId, "metal-cli-events-dev")
				if err != nil {
					t.Error(err)
				}
				status, err = helper.IsDeviceStateActive(deviceId, "active")
				if err != nil {
					status, err = helper.IsDeviceStateActive(deviceId, "active")
					if err != nil || !status {
						t.Error(err)
					}
				}
				root.SetArgs([]string{subCommand, "get", "-d", deviceId})
				rescueStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w
				if err := root.Execute(); err != nil {
					t.Error(err)
				}
				w.Close()
				out, _ := io.ReadAll(r)
				os.Stdout = rescueStdout
				if !strings.Contains(string(out[:]), "Queued for provisioning") &&
					!strings.Contains(string(out[:]), "Connected to magic install system") &&
					!strings.Contains(string(out[:]), "Provision complete! Your device is ready to go.") {
					t.Error("expected output should include Queued for provisioning in output string")
				}
				err = helper.CleanTestDevice(deviceId)
				if err != nil {
					t.Error(err)
				}
				err = helper.CleanTestProject(projectId)
				if err != nil {
					t.Error(err)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := rootClient.NewCommand()
			rootCmd.AddCommand(tt.fields.MainCmd)
			tt.cmdFunc(t, tt.fields.MainCmd)
		})
	}
}
