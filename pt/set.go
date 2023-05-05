package pt

import (
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/promptc/promptc-go/prompt"
)

type PromptSet []*prompt.PromptC

func NewSet(paths []string) PromptSet {
	var set []*prompt.PromptC
	for _, path := range paths {
		txt, err := iox.ReadAllText(path)
		if err != nil {
			panic(err)
		}
		pt := prompt.ParsePromptC(txt)
		set = append(set, pt)
	}
	return set
}

func (set PromptSet) Get(ability string) *prompt.PromptC {
	if len(set) == 0 {
		return nil
	}
	for _, pt := range set {
		if pt.Project == ability {
			return pt
		}
	}
	return nil
}
