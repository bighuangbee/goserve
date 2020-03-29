package user

import (
	"github.com/bighuangbee/gogin/model"
	"github.com/bighuangbee/gogin/pkg/message"
	"github.com/bighuangbee/gogin/pkg/respone"
	"github.com/bighuangbee/gogin/pkg/unit"
	"github.com/bighuangbee/gogin/service/userService"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func Login(c *gin.Context){

	type UsersValid struct {
		Username	string `form:"username" binding:"required"`
		Password	string `form:"password" binding:"required"`
	}

	var usersValid UsersValid
	err := c.ShouldBind(&usersValid)
	if err != nil {
		respone.Error(c, respone.INVALID_PARAMS, err.Error())
		return
	}

	code, data := userService.LoginCheck(usersValid.Username, usersValid.Password)
	if(code == respone.SUCCESS){
		respone.Success(c, message.USER_LOGIN_SUCCESS, data)
	}else{
		respone.Error(c, code, "")
	}
}

func Logout(c *gin.Context){
	userService.Logout(c.GetInt("user_id"))

	respone.Success(c, message.USER_LOGOUT_SUCCESS, nil)
}

/*
	获取用户列表
 */
func List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	length, _ := strconv.Atoi(c.DefaultQuery("length", "10"))

	userList := model.GetUserList(page, length)

	list := make([]map[string]interface{}, 0)
	for _, user := range userList{
		tmp := make(map[string]interface{})
		tmp["id"] = user.Id
		tmp["username"] = user.Username
		tmp["last_login_time"] = user.LastLoginTime
		tmp["create_time"] = user.CreateTime

		roles := make([]map[string]interface{}, 0)
		for _,role := range user.SysRole {
			val := make(map[string]interface{})
			val["id"] = role.Id
			val["name"] = role.Name
			roles = append(roles, val)
		}
		tmp["roles"] = roles

		list = append(list, tmp)
	}

	data := make(map[string]interface{})
	data["total"] = model.GetUserTotal()
	data["list"] = list

	respone.Success(c, message.OK, data)
}

/*
	获取指定用户信息
 */
func Info(c *gin.Context) {

	user_id := c.GetInt("user_id")
	user := model.GetUserWithRole(user_id)

	// 获取用户拥有的角色
	var role_ids []string
	var auth_ids []string
	var roles []string
	for _, item := range user.SysRole {
		role_ids = append(role_ids, strconv.Itoa(item.Id))
	}

	// 获取用户拥有的权限
	for _, item := range model.GetRole(role_ids) {
		a := strings.Split(item.AuthIds, ",")
		auth_ids = append(auth_ids, a...)
		roles = append(roles, item.Name)
	}
	auth_ids = unit.SliceUnique(auth_ids)

	// 获取全部菜单
	allMenu := model.GetMenu(nil)

	// 按权限筛选用户的菜单
	userMenu := []model.SysMenu{}
	for _, menu := range allMenu {
		for _, val := range auth_ids{
			val, _ := strconv.Atoi(val)
			if(menu.Id == val) {
				userMenu = append(userMenu, menu)
			}
		}
	}

	retData := make(map[string]interface{})
	retData["username"] 	= user.Username
	retData["phone"] 		= user.Phone
	retData["head_pic_url"] = user.HeadPicUrl
	retData["sex"] 	= user.Sex
	retData["roles"] = roles
	retData["menu"] = userMenu

	respone.Success(c, message.OK, retData)
}

func Menu(c *gin.Context)  {
	retData := model.GetMenu(nil)
	respone.Success(c, message.OK, retData)
}

/*
	新增用户
 */
func Add(c *gin.Context) {

	type UsersValid struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
		Roles     []string  `form:"roles" binding:"required"`
	}

	user := new(UsersValid)
	if err := c.ShouldBind(user); err != nil {
		respone.Error(c, respone.INVALID_PARAMS, err.Error())
		return
	}

	salt := unit.RandStr(6)
	userModel := model.SysUser{
		Username:		user.Username,
		Password :		userService.CeatePassword(user.Password, salt),
		Salt:			salt,
		LastLoginTime: 	model.Time{},
		CreateTime : 	model.Time(time.Now()),
		CreateUserId: 	c.GetInt("user_id"),
	}

	result := userModel.Add()
	if result.Error != nil{
		 respone.Error(c, respone.USER_CREATED_FAILED, result.Error.Error())
		 return
	}

	model.AddUserRoles(userModel.Id, user.Roles)

	roleModel := model.GetRoleList()
	roles := make([]map[string]interface{}, 0)
	for _,roleId := range user.Roles {
		for _, role := range roleModel{
			if roleId ==  strconv.Itoa(role.Id){
				val := make(map[string]interface{})
				val["id"] = role.Id
				val["name"] = role.Name
				roles = append(roles, val)
			}
		}

	}

	userData := make(map[string]interface{})
	userData["id"] = userModel.Id
	userData["username"] = userModel.Username
	userData["last_login_time"] = userModel.LastLoginTime
	userData["create_time"] = userModel.CreateTime
	userData["roles"] = roles

	respone.Success(c, message.USER_CREATED_SUCCESS, userData)
	return
}

//修改指定用户信息
func Edit(c *gin.Context) {
	type UsersValid struct {
		Id int `form:"id" binding:"required"`
		Username string `form:"username"`
		Roles     []string  `form:"roles"`
		Sex     int  `form:"sex"`
	}

	user := new(UsersValid)
	if err := c.ShouldBind(user); err != nil {
		respone.Error(c, respone.INVALID_PARAMS, err.Error())
		return
	}

	result := model.UpdateUser(user.Id, user)
	if result.Error != nil{
		respone.Error(c, respone.ERROR, result.Error.Error())
		return
	}

	model.UpdateUserRole(user.Id, user.Roles)

	respone.Success(c, message.SUCCESS, nil)
}

//删除用户
func Del(c *gin.Context) {
	type UsersValid struct {
		Id int `form:"id" binding:"required"`
	}

	user := new(UsersValid)
	if err := c.ShouldBind(user); err != nil {
		respone.Error(c, respone.INVALID_PARAMS, err.Error())
		return

	}
	if user.Id == 1 {
		return
	}

	where := make(map[string]interface{})
	where["id"] = user.Id

	updateDate := make(map[string]interface{})
	updateDate["delete_time"] = time.Now().Unix()

	result := model.UpdateUser(user.Id, updateDate)
	if result.Error != nil{
		respone.Error(c, respone.INVALID_PARAMS, result.Error.Error())
		return
	}

	respone.Success(c, message.USER_DELETED_SUCCESS, nil)
	return
}

func UpdatePassword(c *gin.Context){

	type UsersValid struct {
		Old_password	string `form:"old_password" binding:"required"`
		Password	string `form:"password" binding:"required"`
	}

	usersValid := new(UsersValid)
	if err := c.ShouldBind(usersValid); err != nil {
		respone.Error(c, respone.INVALID_PARAMS, err.Error())
		return
	}

	user_id := c.GetInt("user_id")

	userModel := model.GetUser(user_id)

	if userModel.Password != "" && userModel.Password != userService.CeatePassword(usersValid.Old_password, userModel.Salt) {
		respone.Error(c, respone.USER_OLD_PASSWORD_INCORRENT, "")
		return
	}

	data := make(map[string]interface{})
	data["password"] = userService.CeatePassword(usersValid.Password, userModel.Salt)
	model.UpdateUser(user_id, data)

	respone.Success(c, message.USER_UPDATE_PASSWORD_SUCCESS, nil)
	return
}

func Roles(c *gin.Context){
	respone.Success(c, message.OK, model.GetRoleList())
	return
}