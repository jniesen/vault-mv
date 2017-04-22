package main

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"testing"
)

type testSecret struct {
	src  string
	dest string
	data map[string]interface{}
}

func testVaultClient(t *testing.T) *api.Client {
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return client
}

func testPrepareData(t *testing.T, client *api.Client, s testSecret) {
	if _, err := client.Logical().Write(s.src, s.data); err != nil {
		t.Fatalf("err: %s", err)
	}
	//
	//	return func() {
	//		if _, err := client.Logical().Delete(s.dest); err != nil {
	//			t.Fatalf("err: %s", err)
	//		}
	//	}
}

func TestAuthorizedMv(t *testing.T) {
	var tests = map[string]struct {
		source  string
		dest    string
		want    string
		secrets []testSecret
	}{
		"simpleMove": {
			"secret/bar",
			"secret/foo/bar",
			"Moved secret/bar to secret/foo/bar\n",
			[]testSecret{
				{
					src:  "secret/bar",
					dest: "secret/foo/bar",
					data: map[string]interface{}{"value": "bar"},
				},
			},
		},
		"deepMove": {
			"secret/scoop",
			"secret/cone/scoop",
			"Moved secret/scoop/whipcream/cherry to secret/cone/scoop/whipcream/cherry\nMoved secret/scoop/sprinkles to secret/cone/scoop/sprinkles\n",
			[]testSecret{
				{
					src:  "secret/scoop/whipcream/cherry",
					dest: "secret/cone/scoop/whipcream/cherry",
					data: map[string]interface{}{"value": "bar"},
				},
				{
					src:  "secret/scoop/sprinkles",
					dest: "secret/cone/scoop/sprinkles",
					data: map[string]interface{}{"value": "bar"},
				},
			},
		},
	}

	client := testVaultClient(t)
	for k, test := range tests {
		descr := fmt.Sprintf("mv(%q, %q)", test.source, test.dest)

		for _, s := range test.secrets {
			testPrepareData(t, client, s)
		}

		// move the data
		got, err := mv(test.source, test.dest)
		if err != nil {
			t.Errorf("%s: %s failed: %v", k, descr, err)
			continue
		}

		if got != test.want {
			t.Errorf("%s: %s = %q, want %q", k, descr, got, test.want)
		}

		for _, s := range test.secrets {
			// verify the data
			// Move this to a test helper https://youtu.be/yszygk1cpEc
			read, err := client.Logical().Read(s.dest)
			if err != nil {
				t.Fatalf("err: %s", err)
			}

			if fmt.Sprintf("%v", read.Data) != fmt.Sprintf("%v", s.data) {
				t.Errorf("%s: got %s, want %s", k, read.Data, s.data)
			}
		}
	}
}
