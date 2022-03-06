// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/itn/pkg/itn"
	"github.com/spf13/cobra"
)

// TODOs:
//   1. Option to clean-up experiments
//   2. Option to adjust time of rebalance recommendation
//   3. Option to pass an OD instance and have this tool create a matching instance that is spot to test an interruption
//   4. Automated chaos - give this tool a tag or vpc and allow it to randomly interrupt spot instances at will

const (
	version = "version"
)

var versionID string

type Options struct {
	instanceIDs []string
}

func main() {
	options := Options{}
	rootCmd := &cobra.Command{
		Use:   "itn",
		Short: "itn is a simple CLI tool that triggers Amazon EC2 Spot Interruption Termination Notifications (ITNs) and Rebalance Recommendations.",
		Run: func(cmd *cobra.Command, _ []string) {
			if f, _ := cmd.Flags().GetBool(version); f {
				fmt.Println(versionID)
				os.Exit(0)
			}
			ctx := context.Background()
			cfg, err := config.LoadDefaultConfig(ctx)
			if err != nil {
				panic(err)
			}
			if err := itn.New(cfg).Interrupt(context.Background(), options.instanceIDs); err != nil {
				fmt.Printf("❌ %s", err)
				os.Exit(1)
			}
			fmt.Printf("✅ Successfully sent spot rebalance recommendation and ITN to %v\n", options.instanceIDs)
		},
	}
	rootCmd.PersistentFlags().StringSliceVarP(&options.instanceIDs, "instance-ids", "i", []string{}, "instance IDs to interrupt")
	rootCmd.PersistentFlags().BoolP(version, "v", false, "the version")
	rootCmd.Execute()
}