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

	prefix = "env/yo/stage"
	k := strings.Join([]string{prefix, "YO"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("123")}, nil)
	k = strings.Join([]string{prefix, "THING_ID"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("abc123")}, nil)
	k = strings.Join([]string{prefix, "THING_TOKEN"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("yabbadabba")}, nil)

	cases := []struct {
		args  []string
		count int
	}{
		{args: []string{}, count: 3},
		{args: []string{"YOOOOOOOOOOOO"}, count: 0},
		{args: []string{"YO"}, count: 1},
		{args: []string{"YO", "THING_ID"}, count: 2},
		{args: []string{"YO", "THING_ID", "EXTRA"}, count: 2},
		{args: []string{"THING_*"}, count: 2},
		{args: []string{"THING_*", "YO"}, count: 3},
		{args: []string{"*"}, count: 3},
	}

	for _, test := range cases {
		kvs, _ := get(kv, test.args...)
		if len(kvs) != test.count {
			t.Errorf("expected %d but got %d with %v", test.count, len(kvs), test.args)
		}
	}

	prefix = "env"
	kvs, _ := get(kv, []string{}...)
	if len(kvs) != 0 {
		t.Errorf("expected 0 but got %d with not specific enough prefix", len(kvs))
	}
}
