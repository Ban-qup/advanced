package backend

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// SaveVerificationCode сохраняет код верификации в базе данных.
func SaveVerificationCode(ctx context.Context, db *gorm.DB, email, verificationCode string) error {
	// Определяем время истечения срока действия кода
	expirationTime := time.Now().Add(time.Minute * 10)
	// Создаем новый объект VerificationCode
	code := VerificationCode{
		Email:          email,
		Code:           verificationCode,
		ExpirationTime: expirationTime,
	}
	// Добавляем код верификации в базу данных
	if err := db.Create(&code).Error; err != nil {
		return err // Возвращаем ошибку, если не удалось сохранить код верификации
	}
	return nil // Возвращаем nil, если операция выполнена успешно
}

// CheckVerificationCode проверяет соответствие кода верификации в базе данных.
func CheckVerificationCode(ctx context.Context, db *gorm.DB, email, verificationCode string) (bool, error) {
	var code VerificationCode
	// Ищем запись в базе данных по email и коду верификации
	result := db.Where("email = ? AND code = ?", email, verificationCode).First(&code)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil // Возвращаем false и nil, если запись не найдена в базе данных
		}
		return false, result.Error // Возвращаем false и ошибку, если произошла другая ошибка при поиске
	}
	// Проверяем, не истек ли срок действия кода верификации
	if code.ExpirationTime.Before(time.Now()) {
		return false, nil // Возвращаем false и nil, если срок действия кода истек
	}
	// Удаляем запись из базы данных
	if err := db.Delete(&code).Error; err != nil {
		return false, err // Возвращаем false и ошибку, если не удалось удалить запись
	}
	return true, nil // Возвращаем true и nil, если проверка прошла успешно и запись удалена
}
