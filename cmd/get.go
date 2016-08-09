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

	"github.com/deepthawtz/kv/store"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

var (
	namespace string
	deployEnv string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "get ENV key values for a give app/namespace",
	Long:  `get all keys or specific KEY if provided`,
	Run:   Get,
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&namespace, "app", "a", "", "app/namespace to get ENV vars for")
	getCmd.Flags().StringVarP(&deployEnv, "env", "e", "", "environment to get ENV vars for (e.g., stage, production)")
}

// Get fetches key values
func Get(cmd *cobra.Command, args []string) {
	if namespace == "" || deployEnv == "" {
		fmt.Println("must supply --app and --env")
		os.Exit(-1)
	}

	client := store.NewConsulClient()
	kvs, err := get(client, args...)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for _, x := range kvs {
		fmt.Println(x)
	}
}

func get(client *consul.KV, args ...string) ([]*store.KVPair, error) {
	var (
		k       string
		kvpairs consul.KVPairs
		kvpair  *consul.KVPair
		err     error
	)

	if len(args) == 0 {
		fmt.Println("no args, getting all key/values")
		k = strings.Join([]string{"env", namespace, deployEnv}, "/")
		kvpairs, _, err = client.List(k, nil)
	} else {
		fmt.Println("args, getting just key/value for args")
		for _, x := range args {
			k = strings.Join([]string{"env", namespace, deployEnv, x}, "/")
			kvpair, _, err = client.Get(k, nil)
			if kvpair != nil {
				kvpairs = append(kvpairs, kvpair)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	var kvs []*store.KVPair
	for _, kvp := range kvpairs {
		if len(kvp.Value) > 0 {
			v := string(kvp.Value)
			kvs = append(kvs, &store.KVPair{Key: kvp.Key, Value: v})
		}
	}
	fmt.Printf("%d key/values found at %s\n", len(kvs), k)

	return kvs, nil
}
