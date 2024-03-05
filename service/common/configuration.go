package common

import (
	"log/slog"
	"os"
	"encoding/json"
)

type serviceConf struct {
	Port string `json:"port"`
}

type mysqlConf struct {
	Server string `json:"server"`
	Password string `json:"password"`
	User string `json:"user"`
	DBName string `json:"dbName"`
	ConnMaxLifetime int `json:"connMaxLifetime"` 
  	MaxOpenConns int `json:"maxOpenConns"`
  	MaxIdleConns int `json:"maxIdleConns"`
}

type crvConf struct {
	Server string `json:"server"`
  	AppID string `json:"appID"`
	Token string `json:"token"`
}

type Config struct {
	Service serviceConf `json:"service"`
	CRV crvConf `json:"crv"`
	MySQL mysqlConf `json:"mysql"`
}

var gConfig Config

func InitConfig()(*Config){
	slog.Info("init configuation start ...")
	//获取用户账号
	//获取用户角色信息
	//根据角色过滤出功能列表
	fileName := "conf/conf.json"
	filePtr, err := os.Open(fileName)
	if err != nil {
        slog.Error("Open file failed","error",err.Error())
		return nil
    }
    defer filePtr.Close()

	// 创建json解码器
    decoder := json.NewDecoder(filePtr)
    err = decoder.Decode(&gConfig)
	if err != nil {
		slog.Error("json file decode failed.","error",err.Error())
		return nil
	}
	slog.Info("init configuation end")
	return &gConfig
}

func GetConfig()(*Config){
	return &gConfig
}