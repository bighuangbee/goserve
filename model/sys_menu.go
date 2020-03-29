package model

type SysMenu struct {
	Model

	Id      int `json:"id";gorm:"PRIMARY_KEY"`
	Title   string `json:"title"`
	Path	string `json:"path"`
	Pid		int `json:"pid"`
	Icon 	string `json:"icon"`
	Order	int `json:"order"`
	Children []SysMenu `json:"children"`
}

func GetMenu(where map[string]interface{}) (menu []SysMenu) {
	 DB.Model(&SysMenu{}).Where(where).Find(&menu)
	 return
}