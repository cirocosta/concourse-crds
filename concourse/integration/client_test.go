package integration

import (
	"context"
	"testing"

	"github.com/cirocosta/crds/concourse"
)

const (
	url        = "http://localhost:8080"
	user, pass = "test", "test"
	team       = "main"
)

func TestHelloWorld(t *testing.T) {
	ctx := context.Background()

	client, err := concourse.New(ctx, url, user, pass)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 3; i++ {
		err = client.SetPipeline(ctx, team, "test", []byte(`
resources:
- name: test
  type: git
  source: {uri: "test"}
jobs:
- name: test
  plan:
  - get: test`))
		if err != nil {
			t.Fatal(err)
		}
	}

}
