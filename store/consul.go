package store

import (
	"fmt"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// KVPair represents a key/value pair
type KVPair struct {
	Key   string
	Value string
}

func (kvp *KVPair) String() string {
	parts := strings.Split(kvp.Key, "/")
	key := parts[len(parts)-1]
	return fmt.Sprintf("%s=%s", key, kvp.Value)
}

// NewConsulClient returns a consul client
func NewConsulClient() *consul.KV {
	config := consul.Config{
		Address: viper.GetString("consul.host"),
		Scheme:  viper.GetString("consul.scheme"),
	}
	client, err := consul.NewClient(&config)
	if err != nil {
		panic(err)
	}

	return client.KV()
}
