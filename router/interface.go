package router

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

type Streamer interface {
	CreateChatCompletionStream(context.Context, openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error)
}

type Completioner interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error)
}
