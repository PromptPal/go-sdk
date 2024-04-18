package example

import (
	"context"
	"testing"
	"time"

	"github.com/PromptPal/go-sdk/promptpal"
)

const (
	endpoint = "http://localhost:7788"
	token    = "d6e9a6b170784fdfb4ef54417a32f391"
)

func TestExample(t *testing.T) {
	ctx := context.Background()
	// create a client
	oneMinute := 1 * time.Minute
	c := promptpal.NewPromptPalClient(endpoint, token, promptpal.PromptPalClientOptions{
		Timeout: &oneMinute,
	})
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

	if res.PromptID != string(PPPromptEcho) {
		t.Error(res.PromptID)
	}
	if res.ResponseMessage != "hello world" {
		t.Error(res.ResponseMessage)
	}
}
