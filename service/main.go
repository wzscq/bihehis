package main

import (
	"bihehis/common"
	"bihehis/crv"
	"bihehis/registration"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"time"
)

func main() {
	confFile := "conf/conf.json"
	if len(os.Args) > 1 {
		confFile = os.Args[1]
		slog.Info(confFile)
	}
	//初始化配置
	conf := common.InitConfig(confFile)
	if conf == nil {
		slog.Error("init config failed","confFile",confFile)
		return
	}
	//设置log打印文件名和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	//初始化时区
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	registrationController := registration.CreateRegistrationController(conf)

	registrationController.Bind(router)

	router.Run(conf.Service.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
