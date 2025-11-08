package setup

import (
	"simple-blog-system/config"
	"simple-blog-system/config/db"

	"log"
)

// BaseURL base url of api
const BaseURL = "/v1/api"

// CloseDB close connection to db
var CloseDB func() error

type SetupData struct {
	ConfigData  config.Config
	InternalApp InternalAppStruct
}

func Init() SetupData {
	configData := config.GetConfig()

	//DB INIT
	dbConn, err := db.Init(configData.DB.DSN)
	if err != nil {
		log.Println("database error")
	}

	CloseDB = func() error {
		if err := dbConn.CloseConnection(); err != nil {
			return err
		}

		return nil
	}

	internalAppVar := initInternalApp(dbConn.GormDB)

	return SetupData{
		ConfigData:  configData,
		InternalApp: internalAppVar,
	}
}

func initInternalApp(gormDB *db.GormDB) InternalAppStruct {
	var internalAppVar InternalAppStruct

	initAppRepo(gormDB, &internalAppVar)
	initAppService(&internalAppVar)
	initAppHandler(&internalAppVar)

	return internalAppVar
}
