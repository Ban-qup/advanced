package backend

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Инициализация GORM и подключение к PostgreSQL
func dbConnect() (*gorm.DB, error) {
	dsn := "host=localhost user=postgres password=Alikosh12 dbname=final port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

// Создание таблицы пользователей
func autoMigrate() {
	db, err := dbConnect()
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{}, &Chapter{}, &Book{}, &VerificationCode{})
}
