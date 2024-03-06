package registration

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log/slog"
	"strconv"
	"time"
)

type RegNumSrc struct {
	DepartmentID     string `json:"departmentID"`
	OutpatientTypeID string `json:"outpatientTypeID"`
	ID               string `json:"id"`
	Date             string `json:"date"`
	Supply           int    `json:"supply"`
	Used             int    `json:"used"`
}

type Repository interface {
	query(sql string) ([]map[string]interface{}, error)
	Begin() (*sql.Tx, error)
	GetRegNumSrc(departmentID, outpatientTypeID string) (*RegNumSrc, error)
	GetRegNumList() (*[]RegNumSrc, error)
	GetCurrentSN() (*[]SNItem, error)
	ConsumeRegNum(tx *sql.Tx, id, userID string) (int64, error)
	CreateRegistration(tx *sql.Tx, regInfo *RegInfo, userID string) (int64, error)
	CreateRegStatusRec(tx *sql.Tx, regInfo *RegInfo, userID string) (int64, error)
}

type DefatultRepository struct {
	DB *sql.DB
}

func (repo *DefatultRepository) CreateRegStatusRec(tx *sql.Tx, regInfo *RegInfo, userID string) (int64, error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sql := "insert into registration_status(registration_id,status,status_type,create_time,update_time,create_user,update_user) values('" +
		regInfo.ID + "','01" + regInfo.RegStatus + "','01','" + nowTime + "','" + nowTime + "','" + userID + "','" + userID + "')"
	slog.Info("CreateRegStatusRec", "sql", sql)
	res, err := tx.Exec(sql)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	return id, nil

}

func (repo *DefatultRepository) CreateRegistration(tx *sql.Tx, regInfo *RegInfo, userID string) (int64, error) {
	nowTime := time.Now().Format("2006-01-02 15:04:05")
	sql := "insert into registration(department_id,outpatient_type_id,patient_id,registration_number,create_time,update_time,create_user,update_user) values('" +
		regInfo.DepartmentID + "','" + regInfo.OutpatientTypeID + "','" +
		regInfo.PatientID + "','" + regInfo.RegNum + "','" + nowTime + "','" + nowTime + "','" + userID + "','" + userID + "')"
	slog.Info("CreateRegistration", "sql", sql)
	res, err := tx.Exec(sql)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (repo *DefatultRepository) ConsumeRegNum(tx *sql.Tx, id, userID string) (int64, error) {
	sql := "update registration_number_source set used=used+1,version=version+1,update_user='" + userID + "' where id='" + id + "' and supply>used"
	res, err := tx.Exec(sql)
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	rowCount, err := res.RowsAffected()
	if err != nil {
		slog.Error(err.Error())
		return 0, err
	}

	return rowCount, nil
}

func (repo *DefatultRepository) GetCurrentSN() (*[]SNItem, error) {
	nowTime := time.Now()
	dateStart := nowTime.Format("2006-01-02") + " 00:00:00"
	dateEnd := nowTime.Format("2006-01-02") + " 23:59:59"
	sql := "select department_id,outpatient_type_id,max(registration_number) as sn from registration where create_time>='" +
		dateStart + "' and create_time<='" + dateEnd + "' group by department_id,outpatient_type_id"
	slog.Info("GetCurrentSN", "sql", sql)
	list, err := repo.query(sql)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	var snItemList []SNItem
	for _, item := range list {
		snItem := SNItem{
			Key:  item["department_id"].(string) + "_" + item["outpatient_type_id"].(string),
			Date: nowTime.Format("2006-01-02"),
		}
		if item["sn"] != nil {
			snItem.SN, _ = strconv.Atoi(item["sn"].(string))
		}
		snItemList = append(snItemList, snItem)
	}

	return &snItemList, nil
}

func (repo *DefatultRepository) GetRegNumList() (*[]RegNumSrc, error) {
	dateStart := time.Now().Format("2006-01-02") + " 00:00:00"
	dateEnd := time.Now().Format("2006-01-02") + " 23:59:59"
	sql := "select id,supply,used,date from registration_number_source where date>='" +
		dateStart + "' and date<='" + dateEnd + "'"

	list, err := repo.query(sql)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	var regNumSrcs []RegNumSrc
	for _, item := range list {
		regNumSrc := RegNumSrc{
			ID:   item["id"].(string),
			Date: item["date"].(string),
		}
		regNumSrc.Supply, _ = strconv.Atoi(item["supply"].(string))
		regNumSrc.Used, _ = strconv.Atoi(item["used"].(string))
		regNumSrcs = append(regNumSrcs, regNumSrc)
	}

	return &regNumSrcs, nil
}

func (repo *DefatultRepository) GetRegNumSrc(departmentID, outpatientTypeID string) (*RegNumSrc, error) {
	dateStart := time.Now().Format("2006-01-02") + " 00:00:00"
	dateEnd := time.Now().Format("2006-01-02") + " 23:59:59"
	sql := "select id,supply,used,date from registration_number_source where department_id='" +
		departmentID + "' and outpatient_type_id='" +
		outpatientTypeID + "' and date>='" +
		dateStart + "' and date<='" + dateEnd + "'"

	list, err := repo.query(sql)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	regNumSrc := &RegNumSrc{
		DepartmentID:     departmentID,
		OutpatientTypeID: outpatientTypeID,
		ID:               list[0]["id"].(string),
		Date:             list[0]["date"].(string),
	}

	regNumSrc.Supply, _ = strconv.Atoi(list[0]["supply"].(string))
	regNumSrc.Used, _ = strconv.Atoi(list[0]["used"].(string))

	return regNumSrc, nil
}

func (repo *DefatultRepository) query(sql string) ([]map[string]interface{}, error) {
	rows, err := repo.DB.Query(sql)
	if err != nil {
		slog.Error(err.Error())
		return nil, nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository) toMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, _ := rows.Columns()
	columns := make([]interface{}, len(cols))
	colPointers := make([]interface{}, len(cols))
	for i, _ := range columns {
		colPointers[i] = &columns[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(colPointers...)
		if err != nil {
			slog.Error(err.Error())
			return nil, nil
		}
		row := make(map[string]interface{})
		for i, colName := range cols {
			val := colPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				row[colName] = string((*val).([]byte))
			default:
				row[colName] = *val
			}
		}
		list = append(list, row)
	}
	return list, nil
}

func (repo *DefatultRepository) Begin() (*sql.Tx, error) {
	return repo.DB.Begin()
}

func (repo *DefatultRepository) Connect(
	server, user, password, dbName string,
	connMaxLifetime, maxOpenConns, maxIdleConns int) {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 user,
		Passwd:               password,
		Net:                  "tcp",
		Addr:                 server,
		DBName:               dbName,
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	repo.DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		slog.Error(err.Error())
	}

	pingErr := repo.DB.Ping()
	if pingErr != nil {
		slog.Error(pingErr.Error())
	}

	repo.DB.SetConnMaxLifetime(time.Minute * time.Duration(connMaxLifetime))
	repo.DB.SetMaxOpenConns(maxOpenConns)
	repo.DB.SetMaxIdleConns(maxIdleConns)
	slog.Info("connect to mysql server " + server)
}
