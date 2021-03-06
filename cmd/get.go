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
	"path"
	"regexp"
	"strings"

	"github.com/deepthawtz/kv/store"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

var (
	prefix string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [KEY]",
	Short: "get key/value pairs for a given key prefix (e.g., env/myapp/stage)",
	Long:  `get all key/value pairs or for specific key if provided`,
	Run:   Get,
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix to get key/values from")
}

// Get fetches key values
func Get(cmd *cobra.Command, args []string) {
	if prefix == "" {
		fmt.Println("must supply key/value path --prefix")
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
	all, _, err := client.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	scoped := scopedKeys(all)
	kvpairs, err := matchingKVPairs(scoped, args)
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

	return kvs, nil
}

func scopedKeys(keys consul.KVPairs) consul.KVPairs {
	// filter just the key/values with matching prefix
	var scoped consul.KVPairs
	for _, kv := range keys {
		if len(kv.Value) > 0 {
			p := len(strings.Split(prefix, "/"))
			k := len(strings.Split(string(kv.Key), "/"))
			if p == k-1 {
				scoped = append(scoped, kv)
			}
		}
	}
	return scoped
}

// partitionArgs filters args and returns both args containing wildcards
// and args without wildcards (specific)
func partitionArgs(args []string) (wildcards []string, specific []string) {
	for _, x := range args {
		if strings.Contains(x, "*") {
			x = strings.Replace(x, "*", ".*", -1)
			wildcards = append(wildcards, x)
		} else {
			specific = append(specific, x)
		}
	}

	return wildcards, specific
}

func matchingKVPairs(scoped consul.KVPairs, args []string) (consul.KVPairs, error) {
	var kvpairs consul.KVPairs
	if len(args) == 0 {
		return scoped, nil
	}

	wildcards, specific := partitionArgs(args)
	for _, wc := range wildcards {
		for _, k := range scoped {
			m, err := regexp.MatchString(wc, k.Key)
			if err != nil {
				return nil, err
			}
			if m {
				kvpairs = append(kvpairs, k)
			}
		}
	}

	for _, s := range specific {
		for _, k := range scoped {
			key := path.Join(prefix, s)
			if key == string(k.Key) {
				kvpairs = append(kvpairs, k)
			}
		}
	}

	return kvpairs, nil
}
