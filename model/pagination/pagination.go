package pagination

import (
	"github.com/jinzhu/gorm"
	"goserve/pkg/dbModel"
)

type Page struct {
	Page 	int `form:"page"`
	Length 	int `form:"length"`
}

type Result struct {
	List 	interface{}
	Total 	int64
}

/**
	分页接口
 */
type Pagination interface {
	Pagination(Page)
}

/*
	根据分页接口的实现对象，动态创建DB分页模型
	@param pagination 接收实现了分页接口的模型（对象）作为参数
*/
func Model(pagination Pagination, page Page) (db *gorm.DB){

	if page.Page <= 0 {
		page.Page = 1
	}
	if page.Length <= 0 {
		page.Length = 10
	}

	db = dbModel.DB.Model(pagination).Offset((page.Page-1)* page.Length).Limit(page.Length)
	return
}