package backend

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword хеширует указанный пароль с помощью bcrypt.
// Параметры:
// - password: пароль для хеширования.
// Возвращает хешированный пароль и ошибку, если хеширование не удалось.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash проверяет, соответствует ли указанный пароль хешу.
// Параметры:
// - password: пароль для проверки.
// - hash: хеш, с которым сравнивается пароль.
// Возвращает true, если пароль соответствует хешу, и false в противном случае.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
