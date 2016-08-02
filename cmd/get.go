// Copyright © 2016 Dylan Clendenin
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
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	namespace string
	deployEnv string
)

// KVPair represents a Consul key/value pair
type KVPair struct {
	Key   string
	Value string
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get ENV key values for a give app/namespace",
	Long:  ``,
	Run:   get,
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&namespace, "app", "a", "", "app/namespace to get ENV vars for")
	getCmd.Flags().StringVarP(&deployEnv, "env", "e", "", "environment to get ENV vars for (e.g., stage, production)")
}

func get(cmd *cobra.Command, args []string) {
	if namespace == "" || deployEnv == "" {
		fmt.Println("must supply --app and --env")
		os.Exit(-1)
	}

	client := NewConsulClient()
	kv := client.KV()
	kvpairs, _, err := kv.List(strings.Join([]string{"/env", namespace, deployEnv}, "/"), nil)
	if err != nil {
		panic(err)
	}

	var kvs []KVPair
	for _, kvp := range kvpairs {
		kvs = append(kvs, KVPair{Key: kvp.Key, Value: string(kvp.Value)})
	}
	for _, x := range kvs {
		parts := strings.Split(x.Key, "/")
		key := parts[len(parts)-1]
		if key != "" {
			fmt.Printf("%s=%s\n", key, x.Value)
		}
	}
}

// NewConsulClient returns a consul client
func NewConsulClient() *consul.Client {
	config := consul.Config{
		Address: viper.GetString("consul.host"),
		Scheme:  viper.GetString("consul.scheme"),
	}
	client, err := consul.NewClient(&config)
	if err != nil {
		panic(err)
	}

	return client
}
