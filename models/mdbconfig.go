package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"errors"
)

func InitGorm() {
	if dbOrmDefault != nil {
		return
	}
	db := LinkDBToTest()
	dbOrmDefault = db
}

func CheckDB(db *gorm.DB) bool {
	if db != nil {
		//就问能不能ping通
		errPing := db.DB().Ping()
		if errPing != nil {
			GetLogger().LogErr(errPing, "", "can't connect db")
			return false
		}
		return true
	}
	GetLogger().LogErr(errors.New("db is nil"), "")
	return false
}

func LinkDBToTest() *gorm.DB {
	//db
	db, er := gorm.Open("mysql", MyConfig.dBTestUsername+":"+MyConfig.dBTestPassword+"@tcp("+MyConfig.dBTestHost+")/"+MyConfig.dBTestName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if er != nil {
		//数据库都连不上，启动个毛啊
		GetLogger().LogErr("无法连接测试数据库 "+er.Error(), "wtf_DB_error")
		return nil
	}
	setTheDB(db)

	return db
}

func LinkDBTOProc() *gorm.DB {
	//db
	db, er := gorm.Open("mysql", MyConfig.dBProcUsername+":"+MyConfig.dBProcPassword+"@tcp("+MyConfig.dBProcHost+")/"+MyConfig.dBProcName+"?charset=utf8&parseTime=True&loc=Asia%2FShanghai")
	if er != nil {
		//数据库都连不上，启动个毛啊
		GetLogger().LogErr("无法连接正式数据库 "+er.Error(), "wtf_DB_error")
		return nil
	}

	setTheDB(db)
	return db
}

func setTheDB(db *gorm.DB) {
	db.DB().SetMaxIdleConns(MyConfig.dBMaxIdle)
	db.DB().SetMaxOpenConns(MyConfig.dBMaxConn)

	// 启用Logger，显示详细日志
	db.LogMode(MyConfig.DBlog) // 注释本语句是为了只显示错误日志
	//Docker中时　　将日志打到标准输入输出上 　　其它情况写入日志
	if !isIndocker() {
		logFile, _ := os.OpenFile(MyConfig.DBLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		db.SetLogger(log.New(logFile, "\r\n", log.LstdFlags))
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "t_" + defaultTableName
	}
}
