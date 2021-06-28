package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OneEstimationByID(id uint) (*Estimation, error) {

	est := new(Estimation)

	db, err := open()
	if err != nil {
		return est, err
	}

	err = db.Preload("Groups").Preload("Groups.Items").First(&est, id).Error
	if err != nil {
		return est, err
	}

	return est, nil
}

func SaveEstimation(est *Estimation) error {

	db, err := open()
	if err != nil {
		return err
	}

	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Save(est).Error
	if err != nil {
		return err
	}

	return nil
}

func DeleteObject(obj interface{}) error {
	db, err := open()
	if err != nil {
		return err
	}

	err = db.Delete(obj).Error
	if err != nil {
		return err
	}
	return nil
}

func open() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate() error {
	db, err := open()
	if err != nil {
		return err
	}
	db.AutoMigrate(Estimation{}, Group{}, Item{})
	return nil
}
