package routers

import (
	"github.com/astaxie/beego"
	"ihome_go_2/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{}, "get:GetAreas")
	//处理用户登陆的请求
	beego.Router("/api/v1.0/sessions", &controllers.UserController{}, "post:Login")
	//对房屋首页展示的业务
	beego.Router("/api/v1.0/houses/index", &controllers.HousesIndexController{}, "get:HousesIndex")
	//处理用户session请求
	beego.Router("/api/v1.0/session", &controllers.UserController{}, "get:GetSessionName;delete:DelSessionName")
	//处理用户注册的请求
	beego.Router("/api/v1.0/users", &controllers.UserController{}, "post:Reg")
	//上传文件的请求
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")

	//更新用户名的操作
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:UpdateUserName")
}
