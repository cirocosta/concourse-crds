package concourse_test

import (
	"context"
	"testing"

	"github.com/cirocosta/crds/concourse"
)

func TestRequest(t *testing.T) {
	ctx := context.Background()

	client, err := concourse.New(
		ctx,
		"http://localhost:8080",
		"test",
		"test",
	)
	if err != nil {
		t.FailNow()
	}

	for i := 0; i < 2; i++ {
		err = client.SetPipeline(
			ctx,
			"main",
			"test",
			[]byte(`{
  "resources": [
    {
      "name": "repository",
      "type": "git",
      "source": {
        "uri": "https://github.com/cirocosta/l4"
      }
    }
  ],
  "jobs": [
    {
      "name": "test",
      "plan": [
        {
          "get": "repository"
        }
      ]
    }
  ]
}`),
		)
		if err != nil {
			t.FailNow()
		}
	}
}
