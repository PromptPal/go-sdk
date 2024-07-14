package example

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/PromptPal/go-sdk/promptpal"
)

const (
	endpoint = "http://localhost:7788"
	token    = "2919bb8f00ff44d7822ded033a7c5957"
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

func TestStreamExample(t *testing.T) {
	ctx := context.Background()
	// create a client
	oneMinute := 1 * time.Minute
	c := promptpal.NewPromptPalClient(endpoint, token, promptpal.PromptPalClientOptions{
		Timeout: &oneMinute,
	})
	// call the `Execute` function
	_, err := c.ExecuteStream(
		ctx,
		string(PPPromptEcho),
		PPPromptEchoVariables{
			Text: "hello world",
		},
		nil,
		func(data *promptpal.APIRunPromptResponse) error {
			fmt.Println(data.ResponseMessage)
			return nil
		},
	)

	if err != nil {
		t.Error(err)
	}
}
