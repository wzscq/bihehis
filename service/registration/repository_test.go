package registration

import (
	"bihehis/common"
	"log/slog"
	"testing"
)

func getReporitory() Repository {
	confFile := "../conf/conf.json"
	conf := common.InitConfig(confFile)
	if conf == nil {
		slog.Error("init config failed")
		return nil
	}

	repo := &DefatultRepository{}
	repo.Connect(
		conf.MySQL.Server,
		conf.MySQL.User,
		conf.MySQL.Password,
		conf.MySQL.DBName,
		conf.MySQL.ConnMaxLifetime,
		conf.MySQL.MaxOpenConns,
		conf.MySQL.MaxIdleConns)

	return repo
}

func TestGetRegNumSrc(t *testing.T) {
	repo := getReporitory()
	if repo == nil {
		t.Error("get repository failed")
		return
	}

	regNumSrc, err := repo.GetRegNumSrc("B0000", "01")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if regNumSrc != nil {
		slog.Info("get regNumSrc success", "regNumSrc", regNumSrc)
	} else {
		slog.Info("no registration number source")
	}
}

func TestGetRegNumList(t *testing.T) {
	repo := getReporitory()
	if repo == nil {
		t.Error("get repository failed")
		return
	}

	regNumList, err := repo.GetRegNumList()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if regNumList != nil {
		slog.Info("get regNumList success", "regNumList", regNumList)
	} else {
		slog.Info("no registration number source")
	}
}

func TestGetCurrentSN(t *testing.T) {
	repo := getReporitory()
	if repo == nil {
		t.Error("get repository failed")
		return
	}

	snItemList, err := repo.GetCurrentSN()
	if err != nil {
		t.Error(err.Error())
		return
	}

	if snItemList != nil {
		slog.Info("get snItemList success", "snItemList", snItemList)
	} else {
		slog.Info("no snItemList")
	}
}

func TestConsumeRegNum(t *testing.T) {
	repo := getReporitory()
	if repo == nil {
		t.Error("get repository failed")
		return
	}

	tx, err := repo.Begin()
	if err != nil {
		t.Error(err.Error())
		return
	}

	count, err := repo.ConsumeRegNum(tx, "1", "test")
	if err != nil {
		t.Error(err.Error())
		return
	}

	if count != 1 {
		t.Error("consume regNum failed")
		return
	}

	tx.Rollback()
}

func TestCreateRegistration(t *testing.T) {
	repo := getReporitory()
	if repo == nil {
		t.Error("get repository failed")
		return
	}

	tx, err := repo.Begin()
	if err != nil {
		t.Error(err.Error())
		return
	}

	regInfo := &RegInfo{
		DepartmentID:     "B0000",
		OutpatientTypeID: "01",
		PatientID:        "P0000",
		RegNum:           "0001",
	}

	id, err := repo.CreateRegistration(tx, regInfo, "test")
	if err != nil {
		t.Error(err.Error())
		return
	}

	slog.Info("create registration success", "id", id)

	tx.Commit()
}
