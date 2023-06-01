// Copyright © 2018 Jasmin Gacic <jasmin@stackpointcloud.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package devices

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Client) Retrieve() *cobra.Command {
	retrieveDeviceCmd := &cobra.Command{
		Use:     `get [-p <project_id>] | [-i <device_id>]`,
		Aliases: []string{"list"},
		Short:   "Retrieves device list or device details.",
		Long:    "Retrieves a list of devices in the project, or the details of the specified device. Either a project ID or a device ID is required.",
		Example: `  # Gets the details of the specified device:
  metal device get -i 52b60ca7-1ae2-4875-846b-4e4635223471
  
  # Gets a list of devices in the specified project:
  metal device get -p 5ad070a5-62e8-4cfe-a0b9-3b79e59f1cfe`,

		RunE: func(cmd *cobra.Command, args []string) error {
			deviceID, _ := cmd.Flags().GetString("id")
			projectID, _ := cmd.Flags().GetString("project-id")

			if deviceID == "" && projectID == "" {
				return fmt.Errorf("either id or project-id should be set")
			}
			cmd.SilenceUsage = true
			include := []string{"Inner_example"} // []string | Nested attributes to include. Included objects will return their full attributes. Attribute names can be dotted (up to 3 levels) to included deeply nested objects. (optional)
			exclude := []string{"Inner_example"} // []string | Nested attributes to exclude. Excluded objects will return only the href attribute. Attribute names can be dotted (up to 3 levels) to exclude deeply nested objects. (optional)
			if deviceID != "" {

				device, r, err := c.Service.FindDeviceById(context.Background(), deviceID).Include(include).Exclude(exclude).Execute()
				if err != nil {
					fmt.Printf("Error when calling `DevicesApi.FindDeviceById``: %v\n", err)
					fmt.Printf("Full HTTP response: %v\n", r)
				}
				header := []string{"ID", "Hostname", "OS", "State", "Created"}

				data := make([][]string, 1)
				data[0] = []string{device.GetId(), device.GetHostname(), *device.GetOperatingSystem().Name, device.GetState(), device.GetCreatedAt().String()}

				return c.Out.Output(device, header, &data)
			}
			devices, r, err := c.Service.FindProjectDevices(context.Background(), projectID).Execute()
			if err != nil {
				fmt.Printf("Error when calling `DevicesApi.FindDeviceById``: %v\n", err)
				fmt.Printf("Full HTTP response: %v\n", r)
			}

			data := make([][]string, len(devices.GetDevices()))

			for i, dc := range devices.GetDevices() {
				data[i] = []string{dc.GetId(), dc.GetHostname(), *dc.GetOperatingSystem().Name, dc.GetState(), dc.GetCreatedAt().String()}
			}
			header := []string{"ID", "Hostname", "OS", "State", "Created"}

			return c.Out.Output(devices, header, &data)
		},
	}

	retrieveDeviceCmd.Flags().StringP("project-id", "p", "", "The project's UUID. This flag is required, unless specified in the config created by metal init or set as METAL_PROJECT_ID environment variable.")
	retrieveDeviceCmd.Flags().StringP("id", "i", "", "The UUID of a device.")

	return retrieveDeviceCmd
}
