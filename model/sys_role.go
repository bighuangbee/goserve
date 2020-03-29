package model

type SysRole struct {
	Model

	Id      int `json:"id";gorm:"PRIMARY_KEY"`
	Name    string `json:"name"`
	Order	int `json:"order"`
	Status	int `json:"status"`
	AuthIds string `json:"auth_ids"`
	Remark	string `json:"remark"`
}

func GetRole(ids []string) (roles []*SysRole) {
	 DB.Model(&SysRole{}).Where("id in (?)", ids).Scan(&roles)
	 return
}
func GetRoleList() (roles []*SysRole){
	DB.Model(&SysUser{}).Find(&roles)
	return roles
}