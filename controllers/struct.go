package controllers

import (
	"time"
)

type GwAlert struct {
	Id         int       `json:"Id" xorm:"not null pk autoincr INT(11)"`
	DeviceId   string    `json:"device_id" xorm:"comment('设备ID') VARCHAR(255)"`
	MessageId  string    `json:"message_id" xorm:"not null default '' comment('报警编号') VARCHAR(255)"`
	AlertType  string    `json:"alert_type" xorm:"not null default '10' comment('报警类型:10=压力,20=偷水,30=撞到,40=在线,50=信号强度') ENUM('10','20','30','40','50')"`
	Cola       string    `json:"cola" xorm:"comment('通知参数1') VARCHAR(255)"`
	Colb       string    `json:"colb" xorm:"comment('通知参数2') VARCHAR(255)"`
	Colc       string    `json:"colc" xorm:"comment('通知参数3') VARCHAR(255)"`
	Totala     string    `json:"totala" xorm:"comment('在线') VARCHAR(255)"`
	Totalb     string    `json:"totalb" xorm:"comment('偷水') VARCHAR(255)"`
	Totalc     string    `json:"totalc" xorm:"comment('撞倒') VARCHAR(255)"`
	Totald     string    `json:"totald" xorm:"comment('开机') VARCHAR(255)"`
	Sendtime   int       `json:"sendtime" xorm:"comment('发送时间') VARCHAR(255)"`
	Createtime time.Time `json:"createtime" xorm:"default 'CURRENT_TIMESTAMP' comment('数据插入时间') TIMESTAMP"`
	CompanyId  int       `json:"company_id" xorm:"not null default 0 comment('公司ID') INT(11)"`
	Descrip    string    `json:"descrip" xorm:"comment('备注') VARCHAR(255)"`
}

//用于查看水压的历史记录
type GwPressure struct {
	Id            int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	CompanyId     int    `json:"company_id" xorm:"not null default 0 comment('公司ID') INT(11)"`
	DeviceId      string `json:"device_id" xorm:"comment('设备ID') VARCHAR(255)"`
	Sendtime      int    `json:"sendtime" xorm:"comment('发送时间') VARCHAR(255)"`
	PressureValue string `json:"pressure_value" xorm:"comment('压力值') VARCHAR(255)"`
	MsgId         int    `json:"msg_id" binding:"required"`
}

type GwDevice struct {
	Id        int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Address   string `json:"address" xorm:"not null default '' comment('设备地址') VARCHAR(255)"`
	Lng       string `json:"lng" xorm:"not null default '' comment('经度') VARCHAR(255)"`
	Lat       string `json:"lat" xorm:"not null default '' comment('纬度') VARCHAR(255)"`
	DeviceId  string `json:"device_id" xorm:"not null default '' comment('设备号') VARCHAR(255)" binding:"required"`
	State     string `json:"state" xorm:"not null default '60' comment('当前设备状态:10=压力,20=偷水,30=撞到,40=在线,50=信号强度,60=正常') ENUM('10','20','30','40','50','60')"`
	CompanyId int    `json:"company_id" xorm:"not null default 0 comment('所属公司ID') INT(11)" binding:"required"`
	Status    int    `json:"status" xorm:"not null default 0 comment('设备是否安装') INT(11)"`
	AlertId   int    `json:"alert_id" xorm:"comment('报警ID') INT(11)"`
	Signal    string `json:"signal" xorm:"comment('信号值') VARCHAR(255)"`
	Beattime  int    `json:"beattime"`
}

type GwCompany struct {
	Id         int    `json:"Id" xorm:"not null pk autoincr INT(11)"`
	Name       string `json:"name" xorm:"not null default '' comment('公司名称') VARCHAR(255)"`
	Address    string `json:"address" xorm:"default '' comment('公司地址') VARCHAR(255)"`
	Value1     string `json:"value1" xorm:"not null default '0.2' comment('压力阀值1') VARCHAR(255)"`
	Value2     string `json:"value2" xorm:"not null default '0.35' comment('压力阀值2') VARCHAR(255)"`
	Createtime string `json:"createtime" xorm:"-"`
	Email      string `json:"email" xorm:"not null default '' VARCHAR(255)"`
	Tel        string `json:"tel" xorm:"not null default '' VARCHAR(255)"`
	Manager    string `json:"manager" xorm:"not null default '' VARCHAR(255)"`
}

type AlertRequest struct {
	AlertType   int    `json:"type" binding:"required"`
	DeviceId    string `json:"device_id" binding:"required"`
	NoteCnt     int    `json:"note_cnt" binding:"required"`
	WaterCnt    int    `json:"water_cnt" binding:"required"`
	BreakCnt    int    `json:"break_cnt" binding:"required"`
	PoweronTime int    `json:"poweron_time" binding:"required"`
}

type MsgIdTime struct {
	MsgId    int `json:"msg_id" binding:"required"`
	SendTime int `json:"time" binding:"required"`
}

//压力

type WaterPressure struct {
	AlertRequest
	Pressure []WaterPressureValue `json:"pressure" binding:"required"`
}
type WaterPressureValue struct {
	Value    float64 `json:"value" binding:"required"`  //压力值kg/cm^2
	MsgId    int     `json:"msg_id" binding:"required"` //信息编号
	SendTime int     `json:"time" binding:"required"`   //信息时间戳
}

//偷水
type StealWater struct {
	AlertRequest
	WaterStatus float64 `json:"water_status" binding:"required"` //阀门当前开启状态:0->恢复, 非0->当前开启圈数
	WaterMax    float64 `json:"water_max"`                       //阀门开启最大圈数
	WaterTime   float64 `json:"water_time`
	MsgIdTime
}

//撞到
type Bump struct {
	AlertRequest
	BreakStatus int `json:"break_status"`
	BreakTime   int `json:"break_time"`
	MsgIdTime
}

//在线通知
type OnlineNotify struct {
	AlertRequest
	MsgIdTime
}

//信号强度
type SignalStrength struct {
	AlertRequest
	Signal int `json"signal" binding:"required"`
	MsgIdTime
}
