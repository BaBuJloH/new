package server

import (
	"net/http"
	"new/server/order"
	"new/service"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Port string
}

func New(port string, srvc *service.Service) *gin.Engine {
	g := gin.Default()

	g.GET("/about", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNoContent, nil)
	})

	{
		orders := order.Handlers{*srvc}
		g.GET("/orders", orders.List)
		g.GET("/order/:order_uid", orders.Get)
		g.POST("/order", orders.Add)
		g.DELETE("/order/:order_uid", orders.Del)
	}

	return g
}
