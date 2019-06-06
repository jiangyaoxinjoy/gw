package controllers

import (
	"encoding/json"
	"gw2/config"
	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

func (api *BaseController) Get(c *gin.Context) {
	var (
		alert           GwAlert
		device          GwDevice
		company         GwCompany
		pressureHistory []GwPressure
		alertRequest    AlertRequest
		waterPressure   WaterPressure
		stealWater      StealWater
		bump            Bump
		onlineNotify    OnlineNotify
		signalStrength  SignalStrength
	)
	rawData, _ := c.GetRawData()
	json.Unmarshal(rawData, &alertRequest)
	// if err := c.ShouldBindJSON(&alertRequest); err != nil {
	// 	c.JSON(200, gin.H{"Status": -1, "Msg": fmt.Sprintf("json err: %v", err.Error())})
	// 	return
	// }
	db, _ := Connect()

	//TODO GET COMPANY ID
	if found, _ := db.Where("device_id = ?", alertRequest.DeviceId).Get(&device); found == false {
		c.JSON(200, gin.H{"Status": -1, "Msg": "no this device "})
		return
	}
	companyId := device.CompanyId

	alert.CompanyId = companyId
	alert.DeviceId = alertRequest.DeviceId
	alert.AlertType = strconv.Itoa(alertRequest.AlertType)
	alert.Totala = strconv.Itoa(alertRequest.NoteCnt)
	alert.Totalb = strconv.Itoa(alertRequest.WaterCnt)
	alert.Totalc = strconv.Itoa(alertRequest.BreakCnt)
	alert.Totald = strconv.Itoa(alertRequest.PoweronTime)

	if alertRequest.AlertType == 10 {
		json.Unmarshal(rawData, &waterPressure)
		alert.MessageId = strconv.Itoa(waterPressure.Pressure[0].MsgId)
		alert.Sendtime = waterPressure.Pressure[0].SendTime
		pressureHistory = make([]GwPressure, len(waterPressure.Pressure))
		for i := range waterPressure.Pressure {
			if i == 0 {
				//压力1 为当前的水压
				alert.Cola = strconv.FormatFloat(waterPressure.Pressure[i].Value, 'E', -1, 64)
				//获取公司压力的阀值

				db.Where("id = ?", companyId).Get(&company)

				if alert.Cola < company.Value1 || alert.Cola > company.Value2 {
					//更新水压异常当前设备状态
					device.State = "10"
				} else {
					//更新水压正常当前设备状态
					device.State = "60"
				}
			}
			pressureHistory[i].CompanyId = companyId
			pressureHistory[i].DeviceId = alertRequest.DeviceId
			pressureHistory[i].PressureValue = strconv.FormatFloat(waterPressure.Pressure[i].Value, 'E', -1, 64)
			pressureHistory[i].Sendtime = waterPressure.Pressure[i].SendTime
			pressureHistory[i].MsgId = waterPressure.Pressure[i].MsgId
		}
		//插入水压历史曲线表
		db.Insert(&pressureHistory)
		if _, err := db.Insert(&alert); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}

		device.AlertId = alert.Id
		fmt.Println(device)
		if _, err := db.Update(&device); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}
	}
	if alertRequest.AlertType == 20 {
		json.Unmarshal(rawData, &stealWater)
		alert.MessageId = strconv.Itoa(stealWater.MsgId)
		alert.Sendtime = stealWater.SendTime
		alert.Cola = strconv.FormatFloat(stealWater.WaterStatus, 'E', -1, 64)
		alert.Colb = strconv.FormatFloat(stealWater.WaterMax, 'E', -1, 64)
		alert.Colc = strconv.FormatFloat(stealWater.WaterTime, 'E', -1, 64)
		if _, err := db.Insert(&alert); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": fmt.Sprintf("steal insert err: %v", err.Error())})
			return
		}
		device.AlertId = alert.Id
		//判断是否为0=恢复
		if stealWater.WaterStatus == 0 {
			//更新偷水正常当前设备状态
			device.State = "60"
		} else {
			//更新为偷水异常当前设别状态
			device.State = "20"
		}
		if _, err := db.Update(&device); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}
	}

	if alertRequest.AlertType == 30 {
		json.Unmarshal(rawData, &bump)
		alert.MessageId = strconv.Itoa(bump.MsgId)
		alert.Sendtime = bump.SendTime
		alert.Cola = strconv.Itoa(bump.BreakStatus)
		alert.Colb = strconv.Itoa(bump.BreakTime)

		db.Insert(&alert)
		//判断是否为0=恢复
		if bump.BreakStatus == 0 {
			//更新撞到正常当前设备状态
			device.State = "60"
		} else {
			//更新撞到异常当前设备状态
			device.State = "30"
		}
		device.AlertId = alert.Id
		if _, err := db.Update(&device); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}
	}

	if alertRequest.AlertType == 40 {
		json.Unmarshal(rawData, &onlineNotify)
		alert.MessageId = strconv.Itoa(onlineNotify.MsgId)
		alert.Sendtime = onlineNotify.SendTime

		db.Insert(&alert)
		//更新设备表同步字段
		device.Beattime = alert.Sendtime
		if _, err := db.Update(&device); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}
	}

	if alertRequest.AlertType == 50 {
		json.Unmarshal(rawData, &signalStrength)
		alert.MessageId = strconv.Itoa(signalStrength.MsgId)
		alert.Sendtime = signalStrength.SendTime
		alert.Cola = strconv.Itoa(signalStrength.Signal)
		db.Insert(&alert)
		//更新设备表同步字段
		device.Signal = alert.Cola
		if _, err := db.Update(&device); err != nil {
			c.JSON(200, gin.H{"Status": -1, "Msg": err.Error()})
			return
		}
	}

	c.JSON(200, gin.H{"code": 0, "msg": "Hello World"})
}

func Connect() (*xorm.Engine, error) {
	db, err := xorm.NewEngine("mysql", config.Mysql)
	if err != nil {
		return nil, fmt.Errorf("Mysql Error:" + err.Error())
	}
	return db, err
}
