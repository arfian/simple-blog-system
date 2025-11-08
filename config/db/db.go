package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/go-openapi/strfmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"simple-blog-system/config"
)

type GormDB struct {
	*gorm.DB
}

type dbConfig struct {
	GormDB       *GormDB
	ConnectionDB *sql.DB
}

type AuditLog struct {
	ID            strfmt.UUID4 `json:"id" gorm:"type:uuid;default:uuid_generate_v4()"`
	NameTable     string       `json:"name_table"`
	OperationType string       `json:"operation_type"`
	Query         string       `json:"query"`
	CreatedAt     time.Time    `json:"created_at" gorm:"autoCreateTime"`
}

func (db dbConfig) CloseConnection() error {
	return db.ConnectionDB.Close()
}

func Init(dsn string) (dbConfig, error) {
	var (
		dbConfigVar dbConfig
		loggerGorm  logger.Interface
	)
	configData := config.GetConfig()

	loggerGorm = logger.Default.LogMode(logger.Silent)
	if configData.App.Env == "local" {
		loggerGorm = logger.Default.LogMode(logger.Info)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		Logger:                 loggerGorm,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return dbConfigVar, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return dbConfigVar, err
	}

	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(configData.DB.MaxIdletimeConn))
	sqlDB.SetMaxIdleConns(configData.DB.MaxIdleConn)
	sqlDB.SetMaxOpenConns(configData.DB.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(configData.DB.MaxLifetimeConn))
	dbConfigVar.ConnectionDB = sqlDB

	dbConfigVar.GormDB = &GormDB{gormDB}
	fmt.Print("database connected")
	RegisterCallbacks(gormDB)

	return dbConfigVar, nil
}

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Update().After("gorm:update").Register("update_audit_log", updateAuditLog)
	db.Callback().Create().After("gorm:create").Register("create_audit_log", createAuditLog)
}

func createAuditLog(db *gorm.DB) {
	if db.Statement.Table == "audit_log" {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	logEntry := &AuditLog{
		NameTable:     db.Statement.Schema.Table,
		Query:         sql,
		OperationType: "INSERT",
	}

	logDb := db.Session(&gorm.Session{SkipHooks: true, NewDB: true})
	if err := logDb.Table("audit_log").Save(logEntry).Error; err != nil {
		return
	}
	return
}

func updateAuditLog(db *gorm.DB) {
	if db.Statement.Table == "audit_log" {
		return
	}

	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)

	logEntry := &AuditLog{
		NameTable:     db.Statement.Schema.Table,
		Query:         sql,
		OperationType: "UPDATE",
	}

	logDb := db.Session(&gorm.Session{SkipHooks: true, NewDB: true})
	if err := logDb.Table("audit_log").Save(logEntry).Error; err != nil {
		return
	}
	return
}
