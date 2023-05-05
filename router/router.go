package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/gpt"
	"github.com/promptc/api-server/pt"
	"github.com/promptc/promptc-go/prompt"
	"github.com/promptc/promptc-go/variable/interfaces"
	"io"
	"strings"
)

type Request struct {
	Input map[string]string `json:"input"`
}

func AppendRouter(group *gin.RouterGroup) {
	group.POST("/:ability", AbilityHandler)
	group.POST("/:ability/stream", AbilityStreamHandler)
	group.GET("/:ability/info", AbilityVarHandler)
}

var set pt.PromptSet
var gptS gpt.GptSet

func AbilityHandler(c *gin.Context) {
	req, ptc, err := parseRequest(c)
	if err != nil {
		return
	}
	rst, err := gpt.FeedPrompt(gptS, ptc, req.Input)
	if err != nil {
		c.String(500, "GPT Error")
		return
	}
	c.String(200, rst)
}

func AbilityStreamHandler(c *gin.Context) {
	req, ptc, err := parseRequest(c)
	if err != nil {
		return
	}
	streamHandler(c, ptc, req.Input)
}

func AbilityVarHandler(c *gin.Context) {
	ptc := getAbility(c)
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

var ErrNotFound = errors.New("Not Found")
var ErrBindingError = errors.New("BINDING_ERROR")

func parseRequest(c *gin.Context) (Request, *prompt.PromptC, error) {
	ptc := getAbility(c)
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

func getAbility(c *gin.Context) *prompt.PromptC {
	abilityStr := c.Param("ability")
	abilityStr = strings.TrimSpace(abilityStr)
	if abilityStr == "" {
		return nil
	}
	return set.Get(abilityStr)
}

func streamHandler(c *gin.Context, pt *prompt.PromptC, varMap map[string]string) {
	stream, err := gpt.StreamPrompt(gptS, pt, varMap)
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
