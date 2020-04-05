package model

import (
	"goserve/modules/dbModel"
	"strconv"
)

type SysUserRoles struct {
	Model

	Id        int `json:"id" gorm:"PRIMARY_KEY"`
	UserId    int `json:"user_id" gorm:"column:sys_user_id"`
	RoleId    int `json:"role_id" gorm:"column:sys_role_id"`
}

func AddUserRoles(userId int, data []string) (bool){
	for _, val := range data {
		role_id,_ := strconv.Atoi(val)
		userRole := SysUserRoles{
			UserId: userId,
			RoleId: role_id,
		}
		dbModel.DB.Create(&userRole)
	}
	return true
}

func UpdateUserRole(userId int, data []string) bool{
	dbModel.DB.Where("sys_user_id = ?", userId).Delete(&SysUserRoles{})

	for _, val := range data {
		role_id,_ := strconv.Atoi(val)
		userRole := SysUserRoles{
			UserId: userId,
			RoleId: role_id,
		}
		dbModel.DB.Create(&userRole)
	}
	return true
}
