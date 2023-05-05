package gpt

import (
	"context"
	"fmt"
	scheduler "github.com/promptc/openai-scheduler"
	"github.com/promptc/promptc-go/prompt"
	"github.com/sashabaranov/go-openai"
)

func FeedPrompt(gs scheduler.Scheduler, pt *prompt.PromptC, varMap map[string]string) (string, error) {
	fmt.Println("Start Request:", varMap)
	cli := gs.GetClient()
	compiled := pt.Compile(varMap)
	req := openai.ChatCompletionRequest{
		Model:    pt.GetConf().Model,
		Messages: compiled.OpenAIChatCompletionMessages(true),
	}
	resp, err := cli.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func StreamPrompt(gs scheduler.Scheduler, pt *prompt.PromptC, varMap map[string]string) (*openai.ChatCompletionStream, error) {
	fmt.Println("Start Request:", varMap)
	cli := gs.GetClient()
	compiled := pt.Compile(varMap)
	req := openai.ChatCompletionRequest{
		Model:    pt.GetConf().Model,
		Messages: compiled.OpenAIChatCompletionMessages(true),
	}
	resp, err := cli.CreateChatCompletionStream(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
