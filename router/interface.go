package router

import (
	"context"
	scheduler "github.com/promptc/openai-scheduler"
	"github.com/sashabaranov/go-openai"
)

type Streamer interface {
	CreateChatCompletionStream(context.Context, openai.ChatCompletionRequest) (stream *openai.ChatCompletionStream, err error)
}

type Completioner interface {
	CreateChatCompletion(context.Context, openai.ChatCompletionRequest) (response openai.ChatCompletionResponse, err error)
}

type OpenAIClientProvider interface {
	GetClient() OpenAIClient
}

type OpenAIClient interface {
	Streamer
	Completioner
}

func SchedulerToOpenAIProvider(scheduler *scheduler.Scheduler) OpenAIClientProvider {
	return &provider{client: scheduler}
}

type provider struct {
	client *scheduler.Scheduler
}

func (p *provider) GetClient() OpenAIClient {
	return p.client.GetClient()
}
