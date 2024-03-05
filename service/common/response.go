package common

type CommonRsp struct {
	ErrorCode int `json:"errorCode"`
	Message string `json:"message"`
	Error bool `json:"error"`
	Result interface{} `json:"result"`
	Params map[string]interface{} `json:"params"`
}

type CommonError struct {
	ErrorCode int `json:"errorCode"`
	Params map[string]interface{} `json:"params"`
}

const (
	ResultSuccess = 10000000
	ResultWrongRequest = 10000001
	ResultCreateDirError=10000027
	ResultBase64DecodeError=10000028
	ResultCreateFileError=10000029
	ResultCannotLoginCRV=10200001
	ResultSaveDataError=10200002
	ResultQueryRequestError=10200003
)

var errMsg = map[int]CommonRsp{
	ResultSuccess:CommonRsp{
		ErrorCode:ResultSuccess,
		Message:"操作成功",
		Error:false,
	},
	ResultWrongRequest:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"请求参数错误，请检查参数是否完整，参数格式是否正确",
		Error:true,
	},
	ResultCannotLoginCRV:CommonRsp{
		ErrorCode:ResultWrongRequest,
		Message:"连接到基础数据平台失败，请与管理员联系处理",
		Error:true,
	},
	ResultSaveDataError:CommonRsp{
		ErrorCode:ResultSaveDataError,
		Message:"保存数据到数据时发生错误，请与管理员联系处理",
		Error:true,
	},
	ResultCreateDirError:CommonRsp{
		ErrorCode:ResultCreateDirError,
		Message:"保存文件时创建文件夹失败，请与管理员联系处理",
		Error:true,
	},
	ResultQueryRequestError:CommonRsp{
		ErrorCode:ResultQueryRequestError,
		Message:"调用平台查询数据时发生错误，请与管理员联系处理",
		Error:true,
	},
}

func CreateResponse(err *CommonError,result interface{})(*CommonRsp){
	if err==nil {
		commonRsp:=errMsg[ResultSuccess]
		commonRsp.Result=result
		return &commonRsp
	}

	commonRsp:=errMsg[err.ErrorCode]
	commonRsp.Result=result
	commonRsp.Params=err.Params
	return &commonRsp
}

func CreateError(errorCode int,params map[string]interface{})(*CommonError){
	return &CommonError{
		ErrorCode:errorCode,
		Params:params,
	}
}