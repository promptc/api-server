package gpt

import "github.com/sashabaranov/go-openai"

type GptSet interface {
	GetGPT() *openai.Client
}
