package main

import (
	"github.com/gin-gonic/gin"
	"goserve/pkg/config"
	"goserve/pkg/dbModel"
	"goserve/pkg/loger"
	"goserve/pkg/permissions"
	"goserve/pkg/redis"
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
	//if e.Enforce(sub, obj, act) {
	sub := "alice" 	// 将要访问用户的用户或角色
	obj := "data1" 	// 将要访问的资源
	act := "read" 	// 用户或角色对资源执行的操作
	loger.Info("permissions test", permissions.Enforcer.Enforce(sub, obj, act))

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