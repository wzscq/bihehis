package main

import (
	"bihehis/common"
	"bihehis/crv"
	"bihehis/registration"
	"log"
	"log/slog"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	//设置log打印文件名和行号
  	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowAllOrigins:true,
        AllowHeaders:     []string{"*"},
        ExposeHeaders:    []string{"*"},
        AllowCredentials: true,
    }))


	conf:=common.InitConfig()
	if conf==nil{
		slog.Error("init config failed")
		return
	}

	//crvClinet 用于到crvframeserver的请求
	crvClinet:=crv.CRVClient{
		Server:conf.CRV.Server,
		Token:conf.CRV.Token,
		AppID:conf.CRV.AppID,
	}

	registrationRepo:=&registration.DefatultRepository{}
    registrationRepo.Connect(
        conf.MySQL.Server,
        conf.MySQL.User,
        conf.MySQL.Password,
        conf.MySQL.DBName,
        conf.MySQL.ConnMaxLifetime,
        conf.MySQL.MaxOpenConns,
        conf.MySQL.MaxIdleConns)

	registrationController:=registration.RegistrationController{
		CRVClient:&crvClinet,
		Repository:registrationRepo,
	}

	registrationController.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}