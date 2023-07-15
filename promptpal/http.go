package promptpal

type errorResponse struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"error"`
}

type apiRunPromptPayload struct {
	// it's a struct, not any
	Variables any    `json:"variables"`
	UserId    string `json:"userId"`
}

type APIRunPromptResponse struct {
	PromptID           string `json:"id"`
	ResponseMessage    string `json:"message"`
	ResponseTokenCount int    `json:"tokenCount"`
}
