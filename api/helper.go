package api

import (
	"github.com/gin-gonic/gin"
	"github.com/promptc/promptc-go/prompt"
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
