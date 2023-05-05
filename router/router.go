package router

import (
	"github.com/gin-gonic/gin"
)

func (p *Provider) AppendRouter(group *gin.RouterGroup) {
	group.POST("/:ability", p.AbilityHandler)
	group.POST("/:ability/stream", p.AbilityStreamHandler)
	group.GET("/:ability/info", p.AbilityVarHandler)
}
