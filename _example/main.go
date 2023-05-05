package main

import (
	"encoding/json"
	"github.com/KevinZonda/GoX/pkg/iox"
	"github.com/gin-gonic/gin"
	"github.com/promptc/api-server/api"
	"github.com/promptc/api-server/cast"
	scheduler "github.com/promptc/openai-scheduler"
)

func main() {

	// Load config
	cfgStr := e(iox.ReadAllText("config.json"))
	var cfg cfgModel
	if _e := json.Unmarshal([]byte(cfgStr), &cfg); _e != nil {
		panic(_e)
	}

	// Init Prompt API Provider
	openaiScheduler := scheduler.NewScheduler(cfg.Tokens)
	apiProvider := api.NewProvider(
		cast.SchedulerToOpenAIProvider(openaiScheduler),
		[]string{
			"echo.promptc",
			"hello.promptc",
		},
	)

	// Init Gin HTTP Server
	engine := gin.New()
	group := engine.Group("/v1")

	// Append API Router
	// /v1/echo -> echo.promptc
	// /v1/hi   -> hello.promptc (project: hi)
	// /v1/hi/stream   -> hello.promptc (project: hi, stream: true)
	// /v1/echo/stream -> echo.promptc  (stream: true)
	// /v1/echo/info   -> echo.promptc
	// /v1/hi/info     -> hello.promptc (project: hi)
	apiProvider.AppendRouter(group)

	// Start HTTP Server
	err := engine.Run(":8080")
	if err != nil {
		panic(err)
	}
}

type cfgModel struct {
	Tokens []string `json:"tokens"`
}

func e[T any](i T, err error) T {
	if err != nil {
		panic(err)
	}
	return i
}
