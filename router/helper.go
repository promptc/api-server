package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/gpt"
	"github.com/promptc/promptc-go/prompt"
	"io"
	"strings"
)

func (p *Provider) getAbility(c *gin.Context) *prompt.PromptC {
	abilityStr := c.Param("ability")
	abilityStr = strings.TrimSpace(abilityStr)
	if abilityStr == "" {
		return nil
	}
	return p.PromptSet.Get(abilityStr)
}

func (p *Provider) streamHandler(c *gin.Context, pt *prompt.PromptC, varMap map[string]string) {
	stream, err := gpt.StreamPrompt(p.Scheduler, pt, varMap)
	if err != nil {
		c.String(500, "GPT Error")
		return
	}
	if stream == nil {
		c.String(500, "GPT Not Available")
		return
	}
	c.Stream(func(w io.Writer) bool {
		r, _err := stream.Recv()
		if _err != nil {
			if errors.Is(_err, io.EOF) {
				c.Writer.WriteHeader(200)
				return false
			}
			c.Writer.WriteHeader(500)
			w.Write([]byte("Something Happened!"))
			return false
		}
		content := r.Choices[0].Delta.Content
		if content == "" {
			return true
		}
		w.Write([]byte(content))
		return true
	})
}

func (p *Provider) parseRequest(c *gin.Context) (Request, *prompt.PromptC, error) {
	ptc := p.getAbility(c)
	if ptc == nil {
		c.String(404, ErrNotFound.Error())
		return Request{}, nil, ErrNotFound
	}
	var req Request
	err := c.BindJSON(&req)
	if err != nil {
		c.String(500, ErrBindingError.Error())
		return Request{}, nil, ErrBindingError
	}
	return req, ptc, nil
}
