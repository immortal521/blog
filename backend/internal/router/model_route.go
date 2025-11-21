package router

import (
	"blog-server/internal/handler"

	"github.com/labstack/echo/v4"
)

func RegisterModelRoutes(r *echo.Group, h handler.IModelHandler) {
	group := r.Group("/model")
	group.POST("/summary", h.ProcessContent)
	group.GET("/sse", h.SSE)
}
