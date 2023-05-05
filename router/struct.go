package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/gpt"
	"github.com/promptc/api-server/pt"
	scheduler "github.com/promptc/openai-scheduler"
	"github.com/promptc/promptc-go/prompt"
	"github.com/promptc/promptc-go/variable/interfaces"
	"io"
)

type Provider struct {
	CliProvider OpenAIClientProvider
	PromptSet   pt.PromptSet
}

func NewProvider(scheduler *scheduler.Scheduler, paths []string) *Provider {
	scheduler.StartDaemon()
	return &Provider{
		CliProvider: scheduler,
		PromptSet:   pt.NewSet(paths),
	}
}

func (p *Provider) AbilityHandler(c *gin.Context) {
	req, ptc, err := p.parseRequest(c)
	if err != nil {
		return
	}
	rst, err := gpt.FeedPrompt(p.CliProvider.GetClient(), ptc, req.Input)
	if err != nil {
		c.String(500, "GPT Error")
		return
	}
	c.String(200, rst)
}

func (p *Provider) AbilityStreamHandler(c *gin.Context) {
	req, ptc, err := p.parseRequest(c)
	if err != nil {
		return
	}
	p.streamHandler(c, ptc, req.Input)
}

func (p *Provider) streamHandler(c *gin.Context, pt *prompt.PromptC, varMap map[string]string) {
	stream, err := gpt.StreamPrompt(p.CliProvider.GetClient(), pt, varMap)
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

func (p *Provider) AbilityVarHandler(c *gin.Context) {
	ptc := p.getAbility(c)
	if ptc == nil {
		c.String(404, ErrNotFound.Error())
		return
	}
	vars := ptc.VarConstraint
	constraints := make(map[string]interfaces.Constraint)
	for _, v := range vars {
		constraints[v.Name()] = v.Constraint()
	}
	c.JSON(200, constraints)
}
