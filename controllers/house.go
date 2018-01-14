package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"ihome_idlefish/models"
	_ "time"
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

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 解析用户数据, 得到房源信息
	user_id := this.GetSession("user_id")
	user := models.User{Id: user_id.(int)}
	house := models.House{User: &user}
	json.Unmarshal(this.Ctx.Input.RequestBody, &house)
	// 2. 插入房源数据到house表中
	o := orm.NewOrm()
	id, err := o.Insert(&house)
	if err != nil {
		resp.Errno = models.RECODE_DATAERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	beego.Info("post houses succ!!! house id =", id)
	// 3. 插入facility和house的多对多关系到表中
	reqMap := make(map[string]interface{})
	json.Unmarshal(this.Ctx.Input.RequestBody, &reqMap)

	facilitys := reqMap["facility"].([]interface{})
	facility := make(models.Facility)
	for _, value := range facilitys {
		facility.Id = value.(int)
	}

	beego.Info("facilitys =", facilitys)
	// 4. 得到新插入house的house_id
	// 5. 返回正确json和house_id
	houseId := make(map[string]interface{})
	houseId["house_id"] = id
	resp.Data = houseId
	return
}

// /api/v1.0/:id([0-9]+)/images [post]
// 上传房源图片信息
func (this *HouseController) UploadImages() {
	beego.Info("=== uploadImages is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
}

// /api/v1.0/user/houses [get]
// 请求当前用户已发布房源信息
func (this *HouseController) GetUserHousesInfo() {
	beego.Info("=== getUserHousesInfo is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 通过seesion得到user_id
	user_id := this.GetSession("user_id")
	// 2. 查询house表 找到所有user_id为当前用户的房屋
	houses := []models.House{}
	o := orm.NewOrm()
	_, err := o.QueryTable("house").Filter("user_id", user_id.(int)).All(&houses)
	if err != nil {
		resp.Errno = models.RECODE_NODATA
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 3. 返回正确json数据
	resp.Data = houses
	return
}

// /api/v1.0/houses/:id([0-9]+) [get]
// 请求房源详细信息
func (this *HouseController) FindHousesById() {
	beego.Info("=== findHousesById is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
}

// /api/v1.0/houses [get]
// 请求房屋搜索信息
func (this *HouseController) GetHousesInfo() {
	beego.Info("=== getHousesInfo is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
}
