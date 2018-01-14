package controllers

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/orm"
	_ "ihome_idlefish/models"
)

type HouseController struct {
	beego.Controller
}

func (this *HouseController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

// /api/v1.0/houses [post]
// 发布房源信息
func (this *HouseController) PostHousesInfo() {
	beego.Info("=== posthouses is called ===")
}

// /api/v1.0/:id([0-9]+)/images [post]
// 上传房源图片信息
func (this *HouseController) UploadImages() {
	beego.Info("=== uploadImages is called ===")
}

// /api/v1.0/user/houses [get]
// 请求当前用户已发布房源信息
func (this *HouseController) GetUserHousesInfo() {
	beego.Info("=== getUserHousesInfo is called ===")
}

// /api/v1.0/houses/:id([0-9]+) [get]
// 请求房源详细信息
func (this *HouseController) FindHousesById() {
	beego.Info("=== findHousesById is called ===")
}

// /api/v1.0/houses [get]
// 请求房屋搜索信息
func (this *HouseController) GetHousesInfo() {
	beego.Info("=== getHousesInfo is called ===")
}
