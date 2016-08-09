package cmd

import (
	"errors"
	"testing"
)

func TestPut(t *testing.T) {
	c, s := makeClient(t)
	defer s.Stop()
	kv := c.KV()

	cases := []struct {
		args   []string
		output error
	}{
		{args: []string{""}, output: errors.New("")},
		{args: []string{"YO"}, output: errors.New("")},
		{args: []string{"YO=123"}, output: nil},
		{args: []string{"YO=123", "THINK_TOKEN=abc123"}, output: nil},
		{args: []string{"YO", "THINK_TOKEN=abc123"}, output: errors.New("")},
	}

	for _, test := range cases {
		if err := set(kv, ""); err == nil {
			t.Errorf("expected %v with arguments (%v) but got: %v", test.output, test.args, err)
		}
	}
}
