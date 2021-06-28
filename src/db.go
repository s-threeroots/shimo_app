package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dsn := os.Getenv("DATABASE_URL")
	connection, err := pq.ParseURL(dsn)
	connection += " sslmode=require"
	sqlDB, err := sql.Open("postgres", connection)
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

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
