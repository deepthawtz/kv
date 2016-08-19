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
	"path"
	"sync"

	"github.com/deepthawtz/kv/store"
	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

// delCmd represents the del command
var delCmd = &cobra.Command{
	Use:   "del KEY [KEY...]",
	Short: "delete key from prefix (e.g., env/myapp/stage)",
	Long:  "",
	Run:   Del,
}

func init() {
	RootCmd.AddCommand(delCmd)

	delCmd.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix to get key/value pairs from")
}

// Del deletes key/value pairs by key
func Del(cmd *cobra.Command, args []string) {
	if prefix == "" {
		fmt.Println("must supply key/value path --prefix")
		os.Exit(-1)
	}

	client := store.NewConsulClient()
	if err := del(client, args...); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func del(client *consul.KV, args ...string) error {
	var wg sync.WaitGroup
	for _, k := range args {
		wg.Add(1)
		go func(k string) {
			key := path.Join(prefix, k)
			_, err := client.Delete(key, nil)
			if err != nil {
				panic(err)
			}
			wg.Done()
		}(k)
	}
	wg.Wait()

	return nil
}
