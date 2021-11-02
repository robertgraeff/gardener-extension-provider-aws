// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"context"
	"fmt"
	"github.com/gardener/gardener-extension-provider-aws/landscaper/pkg/aws"
	"github.com/gardener/landscaper-utils/deployutils/pkg/logger"
	"github.com/gardener/landscaper-utils/deployutils/pkg/utils"
	"github.com/spf13/cobra"
)

// NewDeployCommand creates a *cobra.Command for the deployment of component provider-aws.
func NewDeployCommand(ctx context.Context) *cobra.Command {
	options := utils.NewOptions()

	cmd := &cobra.Command{
		Use:   "landscaper-provider-aws",
		Short: "Deploy component provider-aws. ",
		Long:  "Deploy component provider-aws. ",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := options.Complete(); err != nil {
				fmt.Print(err)
				return err
			}

			return aws.Process(ctx, options)
		},
	}

	logger.InitFlags(cmd.PersistentFlags())
	options.AddFlags(cmd.Flags())
	return cmd
}
