package userService

import (
	"fmt"
	"github.com/bighuangbee/goserve/model"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

const LOGIN_TOKEN_KEY = "sysUser:token_%s";
var PASSWORD_KEY = "e10adc^39A49*ba5&335fX057@2@!#0f853";

type UserClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

/**
	用户登录验证
 */
func LoginCheck(username string, password string) (code int, a map[string]interface{}){
	var data = make(map[string]interface{})

	userModel := model.GetUserByUsername(username)
	if userModel.Username == ""{
		return respone.USER_NOT_EXISTS, data
	}

	//loger.Info(CeatePassword(password, userModel.Salt))

	if userModel.Password != CeatePassword(password, userModel.Salt){
		return respone.USER_PASSWORD_INCORRECT, data
	}

	//登录认证成功后签发jwt token
	claims := UserClaims{
		userModel.Id,
		jwt.StandardClaims{
			ExpiresAt : time.Now().Add(24*7 * time.Hour).Unix(),
			Issuer : "bighuangbee",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtEncrtpy := config.ConfigData.JwtEncrtpy

	tokenStr, err := token.SignedString([]byte(jwtEncrtpy))
	if err == nil{

		data["user_id"] = userModel.Id
		data["token"] = tokenStr

		//保存到redis
		userInfo := make(map[string]interface{})
		userInfo["user_id"] = userModel.Id
		userInfo["username"] = userModel.Username
		userInfo["nickname"] = userModel.Username
		userInfo["token"] = tokenStr

		redis.Redis.HMSet(CreateTokenKey(userModel.Id), userInfo)
		redis.Redis.Expire(CreateTokenKey(userModel.Id), time.Minute * time.Duration(config.ConfigData.LoginExpire))


		updateData := make(map[string]interface{})
		updateData["last_login_time"] = time.Now().Unix()
		model.UpdateUserLoginTime(userModel.Id)

		return respone.SUCCESS, data
	}

	return respone.ERROR, nil
}

func CeatePassword(password string, salt string) (string){
	return unit.MD5(unit.MD5(password) + PASSWORD_KEY + salt)
}

/**
	用户退出登录
 */
func Logout(user_id int ){
	redis.Redis.Del(CreateTokenKey(user_id))
}

/**
 * 生成token的索引键
 * @param $user_id
 * @return string
 */
func CreateTokenKey(user_id int) (string){
	return fmt.Sprintf(LOGIN_TOKEN_KEY, strconv.Itoa(user_id))
}

func GetMenu() []model.SysMenu{
	where := make(map[string]interface{})
	where["pid"] = 0
	menuTree := model.GetMenu(where)

	for key, val := range menuTree {
		where["pid"] = val.Id
		menuTreeSecond := model.GetMenu(where)
		menuTree[key].Children = menuTreeSecond

		for key2, val2 := range menuTreeSecond {
			where["pid"] = val2.Id
			menuTreeSecond[key2].Children = model.GetMenu(where)
		}

	}

	fmt.Println(menuTree)
	return menuTree
}

