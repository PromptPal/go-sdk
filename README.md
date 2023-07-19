# PromptPal Golang SDK

PromptPal is a software application designed to facilitate the collection, storage, modification, and enhancement of prompts.

PromptPal is a versatile service that can be deployed both on-premises and in cloud-native environments, allowing you to set up your own server on your machines.

The Golang SDK is a powerful tool that enables seamless integration with the PromptPal service, providing a convenient way for developers to interact with its features.

## Installation

To install the PromptPal Golang SDK, use the following go get command:

```bash
go get github.com/PromptPal/go-sdk
```

## Usage

#### Prerequirements

* Set up a PromptPal admin panel.
* Download the PromptPal CLI and use the promptpal init and promptpal g commands to generate prompt metadata.

#### SDK Usage

Before using the SDK, make sure to define the endpoint and token. The endpoint refers to the location where your server is deployed (e.g., the service name in a Kubernetes cluster). The token is obtained from the PromptPal admin panel. Please refer to the main repository for more details.

```go
const (
	endpoint = "http://localhost:7788"
	token    = "d6e9a6b170784fdfb4ef54417a32f391"
)
```

```go
import (
	"github.com/PromptPal/go-sdk/promptpal"
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
```
## Used By

This project is utilized by the following projects:

- ClippingKK

## Support

For support, email to annatar.he@gmail.com.

