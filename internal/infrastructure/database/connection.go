package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB(dbURL string) (*gorm.DB, error) {
	gormCfg := &gorm.Config{
		TranslateError: true,
		Logger:         logger.Default.LogMode(logger.Info),
	}
	db, err := gorm.Open(postgres.Open(dbURL), gormCfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}
