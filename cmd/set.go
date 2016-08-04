// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	help = "must supply key/values in form of KEY=VALUE"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set KEY=VAUE [KEY=VALUE...]",
	Short: "set ENV key values for a give app/namespace",
	Long:  `set as many key/value pairs as you wish`,
	Run:   set,
}

func init() {
	RootCmd.AddCommand(setCmd)

	getCmd.Flags().StringVarP(&namespace, "app", "a", "", "app/namespace to get ENV vars for")
	getCmd.Flags().StringVarP(&deployEnv, "env", "e", "", "environment to get ENV vars for (e.g., stage, production)")
}

func set(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println(help)
		os.Exit(-1)
	}

}
