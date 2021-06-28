package main

import (
	"os"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OneEstimationByID(id uint) (*Estimation, error) {

	est := new(Estimation)

	gormDB, err := open()
	if err != nil {
		return est, err
	}

	db, err := gormDB.DB()
	if err != nil {
		return est, err
	}
	defer db.Close()

	err = gormDB.Preload("Groups").Preload("Groups.Items").First(&est, id).Error
	if err != nil {
		return est, err
	}

	return est, nil
}

func SaveEstimation(est *Estimation) error {

	gormDB, err := open()
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = gormDB.Session(&gorm.Session{FullSaveAssociations: true}).Save(est).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteObject(obj interface{}) error {
	gormDB, err := open()
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = gormDB.Delete(obj).Error
	if err != nil {
		return err
	}
	return nil
}

func open() (*gorm.DB, error) {

	dsn := os.Getenv("DATABASE_URL")
	connection, err := pq.ParseURL(dsn)
	connection += " sslmode=require"
	gormDB, err := gorm.Open(postgres.Open(connection), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

func Migrate() error {
	gormDB, err := open()
	if err != nil {
		return err
	}

	db, err := gormDB.DB()
	if err != nil {
		return err
	}
	defer db.Close()

	gormDB.AutoMigrate(Estimation{}, Group{}, Item{})
	return nil
}
