package main

import (
	"github.com/gin-gonic/gin"
	"goserve/pkg/config"
	"goserve/pkg/dbModel"
	"goserve/pkg/loger"
	"goserve/pkg/redis"
	"goserve/routers"
	"net/http"
)

func init(){
	config.Setup("conf/config.yaml")
	loger.Setup("logs")
	dbModel.Setup()
	redis.Setup()
}

func main(){

	gin.SetMode("debug")
	server := &http.Server{
		Addr:           ":" + config.ConfigData.HttpPort,
		Handler:        routers.InitRouter(),
		//ReadTimeout:    10,
		//WriteTimeout:   10,
		MaxHeaderBytes: 1 << 20,
	}

	loger.Info("Server Started Success...", "Listening Port:", config.ConfigData.HttpPort)
	server.ListenAndServe()

}