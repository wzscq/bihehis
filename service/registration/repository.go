package registration

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log/slog"
	"time"
)

type Repository interface {
	query(sql string)([]map[string]interface{},error)
	GetGegistrationNumberSourcen(departmentID,outpatientTypeID string)(string,error)
}

type DefatultRepository struct {
	DB *sql.DB
}

func (repo *DefatultRepository)GetGegistrationNumberSourcen(departmentID,outpatientTypeID string)(string,error){
	dateStart:=time.Now().Format("2006-01-02")+" 00:00:00"
	dateEnd:=time.Now().Format("2006-01-02")+" 23:59:59"
	sql:="select id from registration_number_source where department_id='"+
				departmentID+"' and outpatient_type_id='"+
				outpatientTypeID+"' and date>='"+
				dateStart+"' and date<='"+dateEnd+"'"

	list,err:=repo.query(sql)
	if err!=nil{
		return "",err
	}
	
	if len(list)==0{
		return "",
	}
	return list[0]["number_source_id"].(string),nil
}

func (repo *DefatultRepository)query(sql string)([]map[string]interface{},error){
	rows, err := repo.DB.Query(sql)
	if err != nil {
		slog.Error(err.Error())
		return nil,nil
	}
	defer rows.Close()
	//结果转换为map
	return repo.toMap(rows)
}

func (repo *DefatultRepository)toMap(rows *sql.Rows)([]map[string]interface{},error){
	cols,_:=rows.Columns()
	columns:=make([]interface{},len(cols))
	colPointers:=make([]interface{},len(cols))
	for i,_:=range columns {
		colPointers[i] = &columns[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err:= rows.Scan(colPointers...)
		if err != nil {
			slog.Error(err.Error())	
			return nil,nil
		}
		row:=make(map[string]interface{})
		for i,colName :=range cols {
			val:=colPointers[i].(*interface{})
			switch (*val).(type) {
			case []byte:
				row[colName]=string((*val).([]byte))
			default:
				row[colName]=*val
			} 
		}
		list=append(list,row)
	}
	return list,nil
}

func (repo *DefatultRepository)Connect(
	server,user,password,dbName string,
	connMaxLifetime,maxOpenConns,maxIdleConns int){
	// Capture connection properties.
    cfg := mysql.Config{
        User:   user,
        Passwd: password,
        Net:    "tcp",
        Addr:   server,
        DBName: dbName,
				AllowNativePasswords:true,
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
    slog.Info("connect to mysql server "+server)
}
