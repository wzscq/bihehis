package registration

import (
	"bihehis/crv"
	"bihehis/common"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RegistrationController struct {
	CRVClient *crv.CRVClient
	Repository Repository
}

//Bind bind the controller function to url
func (cl *RegistrationController) Bind(router *gin.Engine) {
	slog.Info("Bind RegistrationController")
	router.POST("registration/register", cl.register)
}

//register create a new registration
//占用号源
//创建一个新的挂号记录
//创建挂号记录对应的状态记录
func (cl *RegistrationController)register(c *gin.Context){
	slog.Info("start registration/register")
	
	
	rsp:=common.CreateResponse(common.CreateError(common.ResultSuccess,nil),nil)
	c.IndentedJSON(http.StatusOK, rsp)
	slog.Info("end registration/register success")
}