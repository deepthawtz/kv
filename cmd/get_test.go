package cmd

import (
	"strings"
	"testing"

	consul "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
)

type configCallback func(c *consul.Config)

func makeClient(t *testing.T) (*consul.Client, *testutil.TestServer) {
	return makeClientWithConfig(t, nil, nil)
}

func makeClientWithConfig(t *testing.T, cb1 configCallback, cb2 testutil.ServerConfigCallback) (*consul.Client, *testutil.TestServer) {
	// Make client config
	conf := consul.DefaultConfig()
	if cb1 != nil {
		cb1(conf)
	}

	// Create server
	server := testutil.NewTestServerConfig(t, cb2)
	conf.Address = server.HTTPAddr

	// Create client
	client, err := consul.NewClient(conf)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	return client, server
}

func TestGet(t *testing.T) {
	c, s := makeClient(t)
	defer s.Stop()
	kv := c.KV()
	kvs, _ := get(kv, "YOOOO")
	if len(kvs) != 0 {
		t.Errorf("expected 0 but got %d", len(kvs))
	}

	namespace = "yo"
	deployEnv = "stage"
	k := strings.Join([]string{"env", namespace, deployEnv, "YO"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("123")}, nil)
	k = strings.Join([]string{"env", namespace, deployEnv, "THING_ID"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("abc123")}, nil)

	kvs, _ = get(kv, "YO")
	if len(kvs) != 1 {
		t.Errorf("expected 1 but got %d", len(kvs))
	}

	kvs, _ = get(kv, "YO", "THING_ID")
	if len(kvs) != 2 {
		t.Errorf("expected 2 but got %d", len(kvs))
	}

	kvs, _ = get(kv, "YO", "THING_ID", "JUNK")
	if len(kvs) != 2 {
		t.Errorf("expected 2 but got %d", len(kvs))
	}
}
