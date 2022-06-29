package database

import (
	transaction "github.com/yktseng/portto-assignment/internal/tx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	conn    *gorm.DB
	txQueue []*transaction.TX
}

func NewDatabase() *Database {
	return &Database{
		txQueue: make([]*transaction.TX, 100),
	}
}

func (d *Database) Connect() error {
	dsn := "host=localhost user=portto password=portto dbname=portto port=5432 sslmode=disable TimeZone=Asia/Taipei"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	d.conn = conn
	sqlDB, err := d.conn.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)
	// d.conn = d.conn
	return nil
}
