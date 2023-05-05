package interfaces

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type OpenAIClientProvider interface {
	GetClient() OpenAIClient
}

type OpenAIClient interface {
	CreateChatCompletionStream(context.Context, openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error)
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error)
}
