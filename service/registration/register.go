package registration

import (
	"errors"
	"log/slog"
)

type RegInfo struct {
	DepartmentID     string `json:"departmentID"`
	OutpatientTypeID string `json:"outpatientTypeID"`
	PatientID        string `json:"patientID"`
	ID               string `json:"id"`
	RegNum           string `json:"regNum"`
	RegStatus        string `json:"regStatus"`
}

/*
Register create a registration record for the patient
regInfo is the registration information
repository is the repository to access the database
*/
func Register(regInfo *RegInfo, repo Repository, snFactory SNFactory, userID string) (*RegInfo, error) {
	regNumSrc, err := repo.GetRegNumSrc(regInfo.DepartmentID, regInfo.OutpatientTypeID)
	if err != nil {
		return nil, err
	}

	if regNumSrc == nil {
		return nil, errors.New("未找到有效的挂号号源")
	}

	if regNumSrc.Supply <= regNumSrc.Used {
		return nil, errors.New("挂号号源已用完")
	}

	//启动事务
	tx, err := repo.Begin()
	if err != nil {
		slog.Error("start transaction failed", "error", err.Error())
		return nil, err
	}
	//占用号源
	rowCount, err := repo.ConsumeRegNum(tx, regNumSrc.ID, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if rowCount != 1 {
		tx.Rollback()
		return nil, errors.New("占用号源失败，请检查号源是否已用完")
	}

	//生成号码
	regInfo.RegNum = snFactory.GetSN(GetRegSNKey(regInfo))
	regInfo.RegStatus = "0"
	//创建挂号记录
	regID, err := repo.CreateRegistration(tx, regInfo, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	regInfo.ID = regID
	//创建挂号状态记录
	_, err = repo.CreateRegStatusRec(tx, regInfo, userID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return regInfo, nil
}

func GetRegSNKey(regInfo *RegInfo) string {
	return regInfo.DepartmentID + "_" + regInfo.OutpatientTypeID
}
