package router

import (
	"goledger-challenge/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	r := gin.Default()

	r.POST("/set", h.SetValue)
	r.GET("/get", h.GetValue)
	r.POST("/sync", h.SyncValue)
	r.GET("/check", h.CheckValue)

	return r
}
