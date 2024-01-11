package gateways

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	root "github.com/equinix/metal-cli/internal/cli"
	"github.com/equinix/metal-cli/internal/gateway"
	outputPkg "github.com/equinix/metal-cli/internal/outputs"
	"github.com/equinix/metal-cli/test/helper"
)

func TestGateways_Create(t *testing.T) {
	subCommand := "gateways"
	rootClient := root.NewClient(helper.ConsumerToken, helper.URL, helper.Version)
	projectName := "metal-cli-" + helper.GenerateRandomString(5) + "-gateways-create"
	project := helper.CreateTestProject(t, projectName)
	vlan := helper.CreateTestVLAN(t, project.GetId())

	tests := []struct {
		name    string
		cmd     *cobra.Command
		want    *cobra.Command
		cmdFunc func(*testing.T, *cobra.Command)
	}{
		{
			name: "create gateways",
			cmd:  gateway.NewClient(rootClient, outputPkg.Outputer(&outputPkg.Standard{})).NewCommand(),
			want: &cobra.Command{},
			cmdFunc: func(t *testing.T, c *cobra.Command) {
				root := c.Root()

				root.SetArgs([]string{subCommand, "create", "-p", project.GetId(), "-v", vlan.GetId(), "-s", "8"})

				out := helper.ExecuteAndCaptureOutput(t, root)

				apiClient := helper.TestClient()
				gateways, _, err := apiClient.MetalGatewaysApi.
					FindMetalGatewaysByProject(context.Background(), project.GetId()).
					Include([]string{"ip_reservation"}).
					Execute()
				if err != nil {
					t.Fatal(err)
				}
				if len(gateways.MetalGateways) != 1 {
					t.Error(errors.New("Gateway Not Found. Failed to create gateway"))
					return
				}

				assertGatewaysCmdOutput(t, string(out[:]), gateways.MetalGateways[0].MetalGateway.GetId(), vlan.GetMetroCode(), strconv.Itoa(int(vlan.GetVxlan())))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := rootClient.NewCommand()
			rootCmd.AddCommand(tt.cmd)
			tt.cmdFunc(t, tt.cmd)
		})
	}
}

func assertGatewaysCmdOutput(t *testing.T, out, gatewayId, metro, vxlan string) {
	if !strings.Contains(out, gatewayId) {
		t.Errorf("cmd output should contain ID of the gateway: [%s] \n output:\n%s", gatewayId, out)
	}

	if !strings.Contains(out, metro) {
		t.Errorf("cmd output should contain metro same as vlan: [%s] \n output:\n%s", metro, out)
	}

	if !strings.Contains(out, vxlan) {
		t.Errorf("cmd output should contain vxlan, gateway is attached with: [%s] \n output:\n%s", vxlan, out)
	}

	if !strings.Contains(out, "ready") {
		t.Errorf("cmd output should contain 'ready' state of the gateway, output:\n%s", out)
	}
}
