package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"ihome_idlefish/models"
	"path"
	"strconv"
	"time"
)

type HouseController struct {
	beego.Controller
}

func (this *HouseController) RetData(resp interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

type HouseInfo struct {
	Area_id    string   `json:"area_id"`    //归属地的区域编号
	Title      string   `json:"title"`      //房屋标题
	Price      string   `json:"price"`      //单价,单位:分
	Address    string   `json:"address"`    //地址
	Room_count string   `json:"room_count"` //房间数目
	Acreage    string   `json:"acreage"`    //房屋总面积
	Unit       string   `json:"unit"`       //房屋单元,如 几室几厅
	Capacity   string   `json:"capacity"`   //房屋容纳的总人数
	Beds       string   `json:"beds"`       //房屋床铺的配置
	Deposit    string   `json:"deposit"`    //押金
	Min_days   string   `json:"min_days"`   //最好入住的天数
	Max_days   string   `json:"max_days"`   //最多入住的天数 0表示不限制
	Facilities []string `json:"facility"`   //房屋设施
}
type RespUserHouses struct {
	Address     string    `json:"address"`
	Area_name   string    `json:"area_name"`
	Ctime       time.Time `json:"ctime"`
	House_id    int       `json:"house_id"`
	Img_url     string    `json:"img_url"`
	Order_count int       `json:"order_count"`
	Price       int       `json:"price"`
	Room_count  int       `json:"room_count"`
	Title       string    `json:"title"`
	User_avatar string    `json:"user_avatar"`
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
	req := HouseInfo{}
	json.Unmarshal(this.Ctx.Input.RequestBody, &req)

	house := models.House{}
	house.Room_count, _ = strconv.Atoi(req.Room_count)
	house.Title = req.Title
	house.Acreage, _ = strconv.Atoi(req.Acreage)
	house.Unit = req.Unit
	house.Deposit, _ = strconv.Atoi(req.Deposit)
	house.Deposit = house.Deposit * 100
	house.Address = req.Address
	house.Price, _ = strconv.Atoi(req.Price)
	house.Price = house.Price * 100
	house.Capacity, _ = strconv.Atoi(req.Capacity)
	house.Beds = req.Beds
	house.Min_days, _ = strconv.Atoi(req.Min_days)
	house.Max_days, _ = strconv.Atoi(req.Max_days)
	area_id, _ := strconv.Atoi(req.Area_id)
	area := models.Area{Id: area_id}
	house.User = &user
	house.Area = &area

	beego.Info("house ======", house)

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
	facilities := []*models.Facility{}

	for _, value := range reqMap["facility"].([]interface{}) {
		fid, _ := strconv.Atoi(value.(string))
		facility := &models.Facility{Id: fid}
		facilities = append(facilities, facility)
	}
	// 建立关系
	m2m := o.QueryM2M(&house, "Facilities")
	_, e := m2m.Add(facilities)
	if e != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

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

	// 1. 得到图片数据
	file, header, err := this.GetFile("house_image")
	if err != nil {
		resp.Errno = models.RECODE_SERVERERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 创建一个file文件的buffer
	fileBuffer := make([]byte, header.Size)
	_, err = file.Read(fileBuffer)
	if err != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 获取后缀
	suffix := path.Ext(header.Filename)

	// 2. 将图片二进制数据存储到fastDFS中, 得到fileID
	groupName, fileId, er := models.FDFSUploadByBuffer(fileBuffer, suffix[1:])
	if er != nil {
		resp.Errno = models.RECODE_IOERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	beego.Info("fdfs upload file succ groupName =", groupName, " fileId =", fileId)

	// 3. 从请求中得到house_id
	house_id := this.Ctx.Input.Param(":id")

	// 4. 查看该房屋中的index_image_url主显图片是否为空 如果为空将Index_image_url设置为此图片的fileID
	house := models.House{}
	house.Id, _ = strconv.Atoi(house_id)

	o := orm.NewOrm()
	if err := o.Read(&house); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 判断index_image_url是否为空
	if house.Index_image_url == "" {
		house.Index_image_url = fileId
	}

	house_image := models.HouseImage{House: &house, Url: fileId}
	house.Images = append(house.Images, &house_image)

	// 5. 将该图片的fileID追加到HouseImage字段中并入库
	if _, err := o.Insert(&house_image); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}
	// 更新house
	if _, err := o.Update(&house); err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	// 6. 拼接完整域名 + fileID路径
	image_url := "http://39.106.152.53/" + fileId
	urlMap := make(map[string]interface{})
	urlMap["url"] = image_url
	// 7. 返回正确json
	resp.Data = urlMap
	return
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
	respHouses := []RespUserHouses{}
	respMap := make(map[string]interface{})
	for _, value := range houses {
		house_temp := RespUserHouses{}
		house_temp.Address = value.Address
		house_temp.Img_url = value.Index_image_url
		house_temp.House_id = value.Id
		house_temp.Ctime = value.Ctime
		house_temp.Order_count = value.Order_count
		house_temp.Price = value.Price
		house_temp.Room_count = value.Room_count
		house_temp.User_avatar = value.User.Avatar_url
		area_tmp := models.Area{}
		o.QueryTable("area").Filter("id", value.Area.Id).One(&area_tmp)
		house_temp.Title = value.Title
		house_temp.Area_name = area_tmp.Name
		user_tmp := models.User{}
		o.QueryTable("user").Filter("id", value.User.Id).One(&user_tmp)
		house_temp.User_avatar = user_tmp.Avatar_url
		respHouses = append(respHouses, house_temp)
	}
	beego.Info("resoHouses =", respHouses)
	respMap["houses"] = respHouses
	resp.Data = respMap
	return
}

// /api/v1.0/houses/:id([0-9]+) [get]
// 请求房源详细信息
func (this *HouseController) FindHousesById() {
	beego.Info("=== findHousesById is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)

	// 1. 从session中获取user_id
	user_id := this.GetSession("user_id")
	// 2. 从url中得到房屋id
	house_id := this.Ctx.Input.Param(":id")
	// 3. 从缓存中取出当前房屋数据 有则直接返回
	cache_conn, err := cache.NewCache("redis", `{"key":"ihome_idlefish", "conn":"127.0.0.1:6379", "dbNum":"0"}`)
	if err != nil {
		resp.Errno = models.RECODE_DBERR
		resp.Errmsg = models.RecodeText(resp.Errno)
		return
	}

	house_info_key := fmt.Sprintf("house_info_%s", house_id)
	house_info_value := cache_conn.Get(house_info_key)
	if house_info_value != nil {
		resData := make(map[string]interface{})
		resData["user_id"] = user_id
		resHouse := make(map[string]interface{})
		json.Unmarshal(house_info_value.([]byte), &resHouse)
		resData["house"] = resHouse
		resp.Data = resData
		return
	}

	beego.Debug("===== no house info desc CACHE!!! SAVE house desc to CACHE ! =====")
	// 4. 缓存中没有 则从数据库中查询
	o := orm.NewOrm()

	house := models.House{}
	house.Id, _ = strconv.Atoi(house_id)
	o.Read(&house)
	o.LoadRelated(&house, "Area")
	o.LoadRelated(&house, "User")
	o.LoadRelated(&house, "Images")
	o.LoadRelated(&house, "Facilities")
	// 5. 关联查询Area, User, Images, Facilities

	// 6. 将房屋详细信息的json格式存入缓存
	house_info_value, _ = json.Marshal(house.To_one_house_desc())
	cache_conn.Put(house_info_key, house_info_value, 3600*time.Second)
	// 7. 返回正确json
	resData := make(map[string]interface{})
	resData["user_id"] = user_id
	resData["house"] = house.To_one_house_desc()

	resp.Data = resData
	return
}

// /api/v1.0/houses [get]
// 请求房屋搜索信息
func (this *HouseController) GetHousesInfo() {
	beego.Info("=== getHousesInfo is called ===")

	resp := Resp{Errno: models.RECODE_OK, Errmsg: models.RecodeText(models.RECODE_OK)}
	defer this.RetData(&resp)
}
