package cmd

import (
	"strings"
	"testing"

	consul "github.com/hashicorp/consul/api"
)

func TestDel(t *testing.T) {
	c, s := makeClient(t)
	defer s.Stop()
	kv := c.KV()

	prefix = "env/yo/stage"
	k := strings.Join([]string{prefix, "YO"}, "/")
	_, _ = kv.Put(&consul.KVPair{Key: k, Value: []byte("123")}, nil)

	if err := del(kv, []string{"YO"}...); err != nil {
		t.Error("expected del to return true but got false")
	}

	kvs, _ := get(kv, []string{"YO"}...)
	if len(kvs) != 0 {
		t.Error("expected key to be deleted")
	}
}
