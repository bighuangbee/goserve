package main

import (
	"github.com/gin-gonic/gin"
	"goserve/modules/config"
	"goserve/modules/dbModel"
	"goserve/modules/loger"
	"goserve/modules/permissions"
	"goserve/modules/redis"
	"goserve/routers"
	"net/http"
)

func init(){
	config.Setup("conf/config.yaml")
	loger.Setup("logs")
	dbModel.Setup()
	redis.Setup()
	permissions.Setup()
}

func main(){

	gin.SetMode("debug")

	server := &http.Server{
		Addr:           ":" + config.ConfigData.HttpPort,
		Handler:        routers.InitRouter(),
		MaxHeaderBytes: 1 << 20,
	}

	loger.Info("Server Started Success...", "Listening Port:", config.ConfigData.HttpPort)
	server.ListenAndServe()

}