package example

import (
	"context"
	"testing"

	"github.com/PromptPal/go-sdk/promptpal"
)

const (
	endpoint = "http://localhost:7788"
	token    = "d6e9a6b170784fdfb4ef54417a32f391"
)

func TestExample(t *testing.T) {
	ctx := context.Background()
	// create a client
	c := promptpal.NewPromptPalClient(endpoint, token)
	// call the `Execute` function
	res, err := c.Execute(
		ctx,
		string(PPPromptEcho),
		PPPromptEchoVariables{
			Text: "hello world",
		},
		nil,
	)

	if err != nil {
		t.Error(err)
	}

	if res != "hello world" {
		t.Error(res)
	}
}
