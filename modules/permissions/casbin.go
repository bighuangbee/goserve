package permissions

import (
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"goserve/modules/config"
	"goserve/modules/dbModel"
	"goserve/modules/loger"
)

var Enforcer *casbin.Enforcer


func Setup(){

	a := gormadapter.NewAdapterByDB(dbModel.DB)
	Enforcer = casbin.NewEnforcer(config.ConfigData.RbacModelFilePath, a)

	// The adapter will use the table named "casbin_rule".
	// If it doesn't exist, the adapter will create it automatically.

	// Load the policy from DB.
	err := Enforcer.LoadPolicy()
	if err != nil {
		loger.Error("Casbin RBAC SetUp Failed ### ", err)
		return
	}

	loger.Info("Casbin SetUp Success...")

	testRbac()
}

func testRbac(){
	/**

	RBAC权限模型

	[p策略] 用户/角色-资源 - 访问规则的映射关系
	p, data2_admin, data2, read
	data2_admin 为用户或角色
	data2 将要访问的资源
	read 对资源进行的操作
	即可表示为： 用户data2_admin对资源data2拥有read的权限 或 角色data2_admin对资源data2拥有read的权限。
	RBAC通常使用第二种方式表示，角色-资源映射关系，然后配合g策略使用

	[g策略] 用户/资源/角色 - 角色的映射关系
	g, alice, data2_admin
	alice可以用来表示用户、资源或角色的其中一种， Cabin 只是将其识别为一个字符串。
	alice属于data2_admin的一个成员。即可表示为：用户alice属于角色data2_admin， 子角色alice属于角色data2_admin

	p	alice1	/user/login	read
	p	alice2	/user/login	read
	p	role1	/user/login	read
	g	user1	role1
	g	user2	role1
	g	user1	role2

	单层RBAC, 不涉及角色继承关系 https://casbin.org/docs/zh-CN/rbac
	多租户扩展 p, admin, domain1, data1, read

	*/
	loger.Info(Enforcer.GetAllRoles())		//所有角色: [role1 role2]
	loger.Info(Enforcer.GetAllSubjects())	//所有策略对象，即包含资源/用户/角色: [alice1 alice2 role1]
	loger.Info(Enforcer.GetAllActions())		//所有操作: [read]
	loger.Info(Enforcer.GetAllObjects())		//所有访问对象（访问资源）: /user/login
	loger.Info(Enforcer.GetUsersForRole("role1"))	//获取属于角色role1的用户: [user1 user2] <nil>
	loger.Info(Enforcer.GetUsersForRole("role2"))	//获取属于角色role1的用户: [user1 user2] <nil>
	loger.Info(Enforcer.GetRolesForUser("user1")) 	//获取用户user1拥有的角色: [role1 role2] <nil>]
	loger.Info(Enforcer.GetRolesForUser("user2")) 	//获取用户user2拥有的角色: [role1] <nil>

	loger.Info(Enforcer.GetPermissionsForUser("role1"))
	aa := Enforcer.GetPermissionsForUser("role1")
	loger.Info(aa[0][1])
	loger.Info(Enforcer.GetGroupingPolicy())

	sub := "role1" 			// 用户或角色
	obj := "/sysUser/login" 	// 将要访问的资源
	act := "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  true]

	sub = "alice1" 			// 用户或角色
	obj = "/sysUser/login" 	// 将要访问的资源
	act = "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  true]

	sub = "alice1" 			// 用户或角色
	obj = "/sysUser/logout" 	// 将要访问的资源
	act = "read" 			// 对资源执行的操作
	loger.Info("permissions Enforce: ", Enforcer.Enforce(sub, obj, act))	// [permissions Enforce：  false]
}