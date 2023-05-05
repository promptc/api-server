package router

import (
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/gpt"
	"github.com/promptc/api-server/pt"
	scheduler "github.com/promptc/openai-scheduler"
	"github.com/promptc/promptc-go/variable/interfaces"
)

type Provider struct {
	Scheduler scheduler.Scheduler
	PromptSet pt.PromptSet
}

func NewProvider(scheduler scheduler.Scheduler, paths []string) *Provider {
	scheduler.StartDaemon()
	return &Provider{
		Scheduler: scheduler,
		PromptSet: pt.NewSet(paths),
	}
}

func (p *Provider) AbilityHandler(c *gin.Context) {
	req, ptc, err := p.parseRequest(c)
	if err != nil {
		return
	}
	rst, err := gpt.FeedPrompt(p.Scheduler, ptc, req.Input)
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
