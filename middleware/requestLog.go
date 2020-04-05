package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"goserve/modules/loger"
	"time"
)

/**
* @Author: bigHuangBee
* @Date: 2020/3/24 19:04
 */


func RequestLog() gin.HandlerFunc {

	return func(c *gin.Context) {

		startTime := time.Now()
		c.DefaultPostForm("test", "")
		c.Request.ParseForm()

		loger.Info(
			"Request " + c.Request.Method,
			c.ClientIP(),
			c.Request.RequestURI,
			c.Request.Header.Get("Authorization"),
			c.Request.PostForm.Encode(),
			c.Request.Header.Get("Content-Type"),
		)

		writer := responeWriter{
			c.Writer,
			bytes.NewBuffer([]byte("")),
		}
		c.Writer = &writer

		c.Next()

		loger.Info(
			"Respone " + c.Request.Method,
			c.ClientIP(),
			c.Request.RequestURI,
			c.Request.Header.Get("Authorization"),
			time.Now().Sub(startTime),
			writer.WriterBuff.String(),
		)
	}
}

/**
	重新实现ResponseWriter接口的Write方法
	保存请求回复的数据副本
 */
type responeWriter struct {
	gin.ResponseWriter
	WriterBuff *bytes.Buffer
}

func (r *responeWriter) Write(body []byte) (size int, err error){
	r.WriterBuff.Write(body)
	size, err = r.ResponseWriter.Write(body)	//调用ResponseWriter接口的原write
	return
}