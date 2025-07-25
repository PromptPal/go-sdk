package promptpal

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type promptPalClient struct {
	endpoint            string
	token               string
	client              *resty.Client
	applyTemporaryToken *func(ctx context.Context) (ApplyTemporaryTokenResult, error)
}

type PromptPalClient interface {
	Execute(ctx context.Context, prompt string, variables any, userId *string) (*APIRunPromptResponse, error)
	ExecuteStream(ctx context.Context, prompt string, variables any, userId *string, onData func(data *APIRunPromptResponse) error) (*APIRunPromptResponse, error)
}

type ApplyTemporaryTokenResult struct {
	Token     string
	Limit     int
	Remaining int
}

type PromptPalClientOptions struct {
	Timeout             *time.Duration
	ApplyTemporaryToken *func(ctx context.Context) (ApplyTemporaryTokenResult, error)
}

func NewPromptPalClient(endpoint string, token string, options PromptPalClientOptions) PromptPalClient {
	client := resty.
		New().
		SetBaseURL(endpoint).
		SetAuthScheme("API").
		SetAuthToken(token).
		SetHeader("User-Agent", "PromptPal-GoSDK/0.1")
		// TODO: collect and report metrics

	if options.Timeout != nil {
		client.SetTimeout(*options.Timeout)
	}

	return &promptPalClient{
		endpoint:            endpoint,
		token:               token,
		client:              client,
		applyTemporaryToken: options.ApplyTemporaryToken,
	}
}

func (p *promptPalClient) ensureTemporaryToken(ctx context.Context) (string, error) {
	if p.applyTemporaryToken == nil {
		return "", nil
	}

	fn := *p.applyTemporaryToken
	result, err := fn(ctx)

	if err != nil {
		return "", err
	}
	return result.Token, nil
}

func (p *promptPalClient) Execute(ctx context.Context, prompt string, variables any, userId *string) (*APIRunPromptResponse, error) {
	payload := apiRunPromptPayload{
		Variables: variables,
	}
	if userId != nil {
		payload.UserId = *userId
	}

	temporaryToken, err := p.ensureTemporaryToken(ctx)
	if err != nil {
		return nil, err
	}

	req := p.client.R().
		SetBody(payload).
		SetPathParam("pid", prompt)

	if temporaryToken != "" {
		req = req.SetHeader(TEMPORARY_TOKEN_HEADER, temporaryToken)
	}

	res, err := req.
		SetResult(APIRunPromptResponse{}).
		SetError(errorResponse{}).
		SetContext(ctx).
		Post("/api/v1/public/prompts/run/{pid}")

	if err != nil {
		return nil, err
	}

	if res.IsError() {
		errMsg := res.Error().(*errorResponse)
		return nil, fmt.Errorf("error: %d %s", errMsg.ErrorCode, errMsg.ErrorMessage)
	}

	result, ok := res.Result().(*APIRunPromptResponse)
	if !ok {
		return nil, errors.New("invalid prompt response type")
	}
	return result, nil
}

func (p *promptPalClient) ExecuteStream(ctx context.Context, prompt string, variables any, userId *string, onData func(data *APIRunPromptResponse) error) (*APIRunPromptResponse, error) {
	payload := apiRunPromptPayload{
		Variables: variables,
	}
	if userId != nil {
		payload.UserId = *userId
	}

	temporaryToken, err := p.ensureTemporaryToken(ctx)
	if err != nil {
		return nil, err
	}

	req := p.client.R().
		SetBody(payload).
		SetPathParam("pid", prompt)

	if temporaryToken != "" {
		req = req.SetHeader(TEMPORARY_TOKEN_HEADER, temporaryToken)
	}

	resp, err := req.
		SetResult(APIRunPromptResponse{}).
		SetError(errorResponse{}).
		SetContext(ctx).
		SetDoNotParseResponse(true).
		Post("/api/v1/public/prompts/run/{pid}/stream")

	if err != nil {
		return nil, err
	}

	defer resp.RawResponse.Body.Close()

	// directly return the response if content has been cached
	if strings.Contains(resp.Header().Get("Content-Type"), "application/json") {
		var value *APIRunPromptResponse
		err := json.NewDecoder(resp.RawResponse.Body).Decode(value)
		if err != nil {
			return nil, err
		}
		return value, nil
	}

	scanner := bufio.NewScanner(resp.RawResponse.Body)
	result := ""
	var lastChunk *APIRunPromptResponse
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		_res := scanner.Text()
		if len(_res) == 0 {
			continue
		}

		if !strings.HasPrefix(_res, "data:") {
			continue
		}

		jsonBuf := []byte(_res[5:])
		var chunkData APIRunPromptResponse
		err = json.Unmarshal(jsonBuf, &chunkData)
		if err != nil {
			return nil, err
		}
		onData(&chunkData)
		result += chunkData.ResponseMessage
		lastChunk = &chunkData
	}

	return &APIRunPromptResponse{
		PromptID:           lastChunk.PromptID,
		ResponseTokenCount: lastChunk.ResponseTokenCount,
		ResponseMessage:    result,
	}, nil
}
