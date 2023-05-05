package router

import (
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/gpt"
	"github.com/promptc/api-server/pt"
	"github.com/promptc/promptc-go/variable/interfaces"
)

type Provider struct {
	GptSet    gpt.GptSet
	PromptSet pt.PromptSet
}

func (p *Provider) AbilityHandler(c *gin.Context) {
	req, ptc, err := p.parseRequest(c)
	if err != nil {
		return
	}
	rst, err := gpt.FeedPrompt(p.GptSet, ptc, req.Input)
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
