package respone

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Json(c *gin.Context, httpCode int, errCode int, data interface{}){
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg: ResponeMsg[errCode],
		Data: data,
	})
}

func Error(c *gin.Context, code int, attachMsg string) {
	msg := ResponeMsg[code]
	if(attachMsg != ""){
		msg = msg + ", " + attachMsg
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code" 	: code,
		"msg" : msg,
		"data"	: nil,
	})
}

func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code" 	: http.StatusOK,
		"msg" : msg,
		"data"	: data,
	})
}

func UnAuthorized(c *gin.Context, msg string){
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"code" 	: http.StatusUnauthorized,
		"msg" : ResponeMsg[USER_AUTHORIZATION] + ", " + msg,
		"data"	: nil,
	})
}