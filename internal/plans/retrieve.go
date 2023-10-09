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

package plans

import (
	"context"
	"fmt"

	metal "github.com/equinix-labs/metal-go/metal/v1"
	"github.com/spf13/cobra"
)

func (c *Client) Retrieve() *cobra.Command {
	return &cobra.Command{
		Use:     `get`,
		Aliases: []string{"list"},
		Short:   "Retrieves a list plans.",
		Long:    "Retrieves a list of plans available to the current user. Response includes plan UUID, slug, and name.",
		Example: `  # Lists the plans available to the current user:
  metal plans get`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			request := c.Service.FindPlans(context.Background()).Include(c.Servicer.Includes(nil)).Exclude(c.Servicer.Excludes(nil))
			filters := c.Servicer.Filters()

			if filters["type"] != "" {
				validType, err := metal.NewFindPlansTypeParameterFromValue(filters["type"])
				if err != nil {
					return err
				}
				request = request.Type_(*validType)
			}

			if filters["slug"] != "" {
				request = request.Slug(filters["slug"])
			}

			plansList, _, err := request.Execute()
			if err != nil {
				return fmt.Errorf("could not list plans: %w", err)
			}

			plans := plansList.GetPlans()
			data := make([][]string, len(plans))

			for i, p := range plans {
				data[i] = []string{p.GetId(), p.GetSlug(), p.GetName()}
			}
			header := []string{"ID", "Slug", "Name"}

			return c.Out.Output(plans, header, &data)
		},
	}
}
