package registration

import (
	"bihehis/common"
	"bihehis/crv"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RegistrationController struct {
	Repository Repository
	SNFactory  SNFactory
}

func CreateRegistrationController(conf *common.Config) *RegistrationController {
	registrationRepo := &DefatultRepository{}
	registrationRepo.Connect(
		conf.MySQL.Server,
		conf.MySQL.User,
		conf.MySQL.Password,
		conf.MySQL.DBName,
		conf.MySQL.ConnMaxLifetime,
		conf.MySQL.MaxOpenConns,
		conf.MySQL.MaxIdleConns)

	snFactory := &DefaultSNFactory{}
	snFactory.Init(registrationRepo)

	return &RegistrationController{
		Repository: registrationRepo,
		SNFactory:  snFactory,
	}
}

// Bind bind the controller function to url
func (cl *RegistrationController) Bind(router *gin.Engine) {
	slog.Info("Bind RegistrationController")
	router.POST("registration/register", cl.register)
}

// register create a new registration
func (cl *RegistrationController) register(c *gin.Context) {
	slog.Info("start registration/register")
	var header crv.CommonHeader
	if err := c.ShouldBindHeader(&header); err != nil {
		slog.Error("register ShouldBindHeader error", "error", err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	var req crv.CommonReq
	if err := c.BindJSON(&req); err != nil {
		slog.Error("register BindJSON error", "error", err)
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, nil), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	if req.List == nil || len(*req.List) == 0 {
		slog.Error("register req.List is nil or empty")
		params := map[string]interface{}{"错误信息": "未提供有效的挂号数据"}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	registerData := (*req.List)[0]
	slog.Info("register", "data", registerData)
	departmentID, ok := registerData["department_id"].(string)
	if !ok {
		slog.Error("register can not get department_id")
		params := map[string]interface{}{"错误信息": "获取科室ID失败"}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	outpatientTypeID, ok := registerData["outpatient_type_id"].(string)
	if !ok {
		slog.Error("register can not get outpatient_type_id")
		params := map[string]interface{}{"错误信息": "获取挂号类型失败"}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	patientID, ok := registerData["patient_id"].(string)
	if !ok {
		slog.Error("register can not get patient_id")
		params := map[string]interface{}{"错误信息": "获取就诊患者ID失败"}
		rsp := common.CreateResponse(common.CreateError(common.ResultWrongRequest, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	regInfo := &RegInfo{
		DepartmentID:     departmentID,
		OutpatientTypeID: outpatientTypeID,
		PatientID:        patientID,
	}

	//执行挂号操作
	regID, err := Register(regInfo, cl.Repository, cl.SNFactory, header.UserID)
	if err != nil {
		slog.Error("register Register error", "error", err)
		params := map[string]interface{}{"错误信息": err.Error()}
		rsp := common.CreateResponse(common.CreateError(common.ResultRegisterError, params), nil)
		c.IndentedJSON(http.StatusOK, rsp)
		return
	}

	slog.Info("register success", "regID", regID)
	rsp := common.CreateResponse(common.CreateError(common.ResultSuccess, nil), nil)
	c.IndentedJSON(http.StatusOK, rsp)
}
