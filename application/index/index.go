package index

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goserve/model"
	"goserve/model/pagination"
	"goserve/pkg/message"
	"goserve/pkg/respone"
)

func Inde(c *gin.Context){
	sysUser := new(model.SysUser).Get(1)

	fmt.Println(sysUser)

	pageParams := pagination.Page{}
	c.BindQuery(&pageParams)

	result := new(model.SysUser).GetList(pageParams)
	//fmt.Println(result)
	respone.Success(c, message.OK, result)
}

func Test(c *gin.Context){
	respone.Success(c, message.OK, "okokoko")
}