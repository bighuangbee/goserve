package routers

import (
	"github.com/gin-gonic/gin"
	"goserve/application/index"
	"goserve/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())       //跨域处理
	r.Use(middleware.RequestLog())
	//r.Use(middleware.Auth())	//处理用户身份鉴权

	r.GET("/", index.Inde)
	indexApi := r.Group("/index/")
	{
		indexApi.POST("test", index.Test)
	}

	return r
}

