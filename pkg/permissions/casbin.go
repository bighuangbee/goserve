package permissions

import (
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"goserve/pkg/config"
	"goserve/pkg/dbModel"
	"goserve/pkg/loger"
)

var Enforcer *casbin.Enforcer

func testRbac(){
	/**
	p策略 -> 用户/角色-资源-访问的映射关系
	p, data2_admin, data2, read
	由于data2_admin可以表示为用户或角色， 映射关系即可以为： 用户data2_admin对资源data2拥有read的权限 或 角色data2_admin对资源data2拥有read的权限。
	RBAC通常使用第二种方式表示，角色-资源映射关系，然后结合g策略

	g策略 -> 用户/资源/角色 - 角色的映射关系
	g, alice, data2_admin
	alice可以是用户、资源或角色的其中一种， Cabin 只是将其识别为一个字符串。
	alice是角色 data2_admin的一个成员。重点关注其中两种表示含义即可：用户alice属于角色data2_admin， 子角色alice属于角色data2_admin

	p	alice1	/user/login	read
	p	alice2	/user/login	read
	p	role1	/user/login	read
	g	user1	role1
	g	user2	role1
	g	user1	role2

	*/
	loger.Info(Enforcer.GetAllRoles())		//所有角色 [role1 role2]
	loger.Info(Enforcer.GetAllSubjects())	//所有策略对象，即包含资源/用户/角色 [alice1 alice2 role1]
	loger.Info(Enforcer.GetAllActions())		//所有操作 [read]
	loger.Info(Enforcer.GetAllObjects())		//所有访问对象（访问资源） /user/login
	loger.Info(Enforcer.GetUsersForRole("role1"))	//获取属于角色role1的用户[user1 user2] <nil>
	loger.Info(Enforcer.GetRolesForUser("user1")) //获取用户user1拥有的角色 [role1 role2] <nil>]
	loger.Info(Enforcer.GetRolesForUser("user2")) //获取用户user2拥有的角色[role1] <nil>

	sub := "role1" 			// 用户或角色
	obj := "/user/login" 	// 将要访问的资源
	act := "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  true]

	sub = "alice1" 			// 用户或角色
	obj = "/user/login" 	// 将要访问的资源
	act = "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  true]

	sub = "alice1" 			// 用户或角色
	obj = "/user/logout" 	// 将要访问的资源
	act = "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  false]
}


func Setup(){

	a := gormadapter.NewAdapterByDB(dbModel.DB)
	Enforcer = casbin.NewEnforcer(config.ConfigData.RbacModelFilePath, a)

	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.

	Enforcer.AddFunction("ParamsMatch", ParamsMatchFunc)

	// Load the policy from DB.
	err := Enforcer.LoadPolicy()
	if err != nil {
		loger.Error("Casbin RBAC SetUp Failed ### ", err)
		return
	}

	testRbac()

	loger.Info("Casbin SetUp Success...")
}

// 自定义规则函数
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	loger.Info("ParamsMatchFunc", args)
	return args, nil
}