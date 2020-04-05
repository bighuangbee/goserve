package pagination

import (
	"github.com/jinzhu/gorm"
	"goserve/modules/dbModel"
)

type Page struct {
	Start  int `form:"start"`
	End int `form:"end"`
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
func Model(pagination Pagination, page Page) (*gorm.DB){

	if page.Start <= 0 {
		page.Start = 1
	}
	if page.End < page.Start {
		page.End = page.Start
	}

	limit := (page.End - page.Start) + 1

	return dbModel.DB.Model(pagination).Offset(page.Start).Limit(limit)
}