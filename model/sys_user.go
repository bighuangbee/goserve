package model

import (
	"goserve/model/pagination"
	"goserve/pkg/dbModel"
)

type SysUser struct {
	Model
	Id            int      `json:"id";gorm:"PRIMARY_KEY"`
	Username      string   `json:"username"`
	Phone         string   `json:"phone"`
	Password      string   `json:"password"`
	Salt          string   `json:"salt"`
	HeadPicUrl    string   `json:"head_pic_url"`
	Sex           int      `json:"sex"`
	CreateUserId  int      `json:"create_user_id"`
	LastLoginTime Time 		`json:"last_login_time"`
	CreateTime    Time 		`json:"create_time"`
	DeleteTime    Time    	`json:"delete_time"`
}

func (model *SysUser) Get(id int) (*SysUser) {
	result := SysUser{}
	dbModel.DB.Model(model).Where("id = ?", id).Scan(&result)
	return &result
}

/*
	实现分页接口, 空接口
*/
func (sysUserModel *SysUser) Pagination(page pagination.Page){
}

func (sysUserModel *SysUser) GetList(page pagination.Page) (*pagination.Result) {

	var sysUser []SysUser
	pagination.Model(sysUserModel, page).Find(&sysUser)

	var count int64
	dbModel.DB.Model(&sysUserModel).Count(&count)

	return &pagination.Result{
		List:  sysUser,
		Total: count,
	}
}