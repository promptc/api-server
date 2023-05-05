package pt

import (
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/promptc/promptc-go/prompt"
)

func LoadPrompt(path string) *prompt.PromptC {
	txt, err := iox.ReadAllText(path)
	if err != nil {
		panic(err)
	}

	f := prompt.ParsePromptC(txt)
	return f
}
