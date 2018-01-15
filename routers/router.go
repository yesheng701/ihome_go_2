package routers

import (
	"github.com/astaxie/beego"
	"ihome_idlefish/controllers"
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
	//请求用户基本信息
	beego.Router("/api/v1.0/user", &controllers.UserController{}, "get:GetUserInfo")
	//上传文件的请求
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{}, "post:UploadAvatar")
	//更新用户名的操作
	beego.Router("/api/v1.0/user/name", &controllers.UserController{}, "put:UpdateUserName")
	//实名认证检查
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{}, "get:GetUserInfo;post:UploadUserAuth")

	/*
	* 房屋相关业务
	 */
	// 发布房源信息
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "post:PostHousesInfo")
	// 上传房源图片信息
	beego.Router("/api/v1.0/houses/:id([0-9]+)/images", &controllers.HouseController{}, "post:UploadImages")
	// 请求当前用户已发布房源信息
	beego.Router("/api/v1.0/user/houses", &controllers.HouseController{}, "get:GetUserHousesInfo")
	// 请求房源详细信息
	beego.Router("/api/v1.0/houses/:id([0-9]+)", &controllers.HouseController{}, "get:FindHousesById")
	// 请求房屋搜索信息
	beego.Router("/api/v1.0/houses", &controllers.HouseController{}, "get:GetHousesInfo")

	/*
	* 订单相关业务
	 */
	//	beego.Router("/api/v1.0/orders", &controller.OrdersController{}, "post:PostOrders")

}
