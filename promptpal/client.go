package promptpal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type promptPalClient struct {
	endpoint string
	token    string
	client   *resty.Client
}

type PromptPalClient interface {
	Execute(ctx context.Context, prompt string, variables any, userId *string) (string, error)
}

func NewPromptPalClient(endpoint string, token string) PromptPalClient {
	client := resty.
		New().
		SetTimeout(10 * time.Second).
		SetBaseURL(endpoint).
		SetAuthScheme("API").
		SetAuthToken(token)

	return &promptPalClient{
		endpoint: endpoint,
		token:    token,
		client:   client,
	}
}

func (p *promptPalClient) Execute(ctx context.Context, prompt string, variables any, userId *string) (string, error) {

	payload := apiRunPromptPayload{
		Variables: variables,
	}
	if userId != nil {
		payload.UserId = *userId
	}

	res, err := p.client.R().
		SetBody(payload).
		SetPathParam("pid", prompt).
		SetResult(apiRunPromptResponse{}).
		SetError(errorResponse{}).
		SetContext(ctx).
		Post("/api/v1/public/prompts/run/{pid}")

	if err != nil {
		return "", err
	}

	if res.IsError() {
		errMsg := res.Error().(*errorResponse)
		return "", fmt.Errorf("error: %d %s", errMsg.ErrorCode, errMsg.ErrorMessage)
	}

	result, ok := res.Result().(*apiRunPromptResponse)
	if !ok {
		return "", errors.New("invalid prompt response type")
	}
	return result.ResponseMessage, nil
}
