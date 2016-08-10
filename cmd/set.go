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
	help = "must supply key/values in form of KEY=VALUE"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set KEY=VAUE [KEY=VALUE...]",
	Short: "set ENV key values for a given key prefix (e.g., env/myapp/stage)",
	Long:  `set as many key/value pairs as you wish`,
	Run:   Set,
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix to get key/value pairs from")
}

// Set sets key/value pairs
func Set(cmd *cobra.Command, args []string) {
	if prefix == "" {
		fmt.Println("must supply --prefix")
		os.Exit(-1)
	}

	if len(args) == 0 {
		fmt.Println(help)
		os.Exit(-1)
	}

	client := store.NewConsulClient()

	if err := set(client, args...); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func set(client *consul.KV, args ...string) error {
	// TODO: use txn when using Consul 0.7
	// var ops []*consul.KVTxnOp
	for _, raw := range args {
		parts := strings.Split(raw, "=")
		if len(parts) != 2 {
			return fmt.Errorf(help)
		}
		k := strings.Join([]string{prefix, parts[0]}, "/")
		v := parts[1]
		fmt.Printf("setting %s = %s\n", k, v)
		if _, err := client.Put(&consul.KVPair{Key: k, Value: []byte(v)}, nil); err != nil {
			return err
		}

		// TODO: use txn when using Consul 0.7
		// ops = append(ops, &consul.KVTxnOp{
		// 	Verb:  "set",
		// Key:  fmt.Sprintf("/env/%s/%s/%s", prefix, k),
		// 	Value: []byte(v),
		// })
	}

	return nil

	// TODO: use txn when using Consul 0.7
	// ops := []*consul.KVTxnOp{
	// 	&consul.KVTxnOp{
	// 		Verb: "get",
	// 		Key:  fmt.Sprintf("/env/%s/%s/%s", prefix, "YO"),
	// 	},
	// }
	// ok, resp, meta, err := kv.Txn(ops, nil)
	// if err != nil {
	// 	fmt.Println(err, resp, meta)
	// 	os.Exit(-1)
	// }
	// if ok {
	// 	fmt.Println(resp)
	// }
}
