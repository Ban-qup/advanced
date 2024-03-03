package backend

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// Структура Claims содержит поля, которые будут включены в JWT токен.
type Claims struct {
	Id                 string `json:"id"`    // Идентификатор пользователя
	Email              string `json:"email"` // Email пользователя
	jwt.StandardClaims        // Стандартные поля JWT токена
	UserType           string `json:"userType"` // Тип пользователя
}

// CreateToken создает новый JWT токен на основе переданных параметров и возвращает его в виде строки.
func CreateToken(id string, email string, userType string) (tokenString string, err error) {
	// Создание структуры Claims с заполненными данными
	claims := &Claims{
		Id:       id,
		Email:    email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Установка времени истечения токена на 24 часа
		},
	}
	// Создание нового JWT токена с указанными claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен с использованием секретного ключа и возвращаем его строковое представление
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", err // Возвращаем ошибку, если не удалось подписать токен
	} else {
		return signedToken, nil // Возвращаем подписанный токен
	}
}

// VerifyToken проверяет переданный JWT токен и возвращает его содержимое.
func VerifyToken(token string) (string, string, string, error) {
	if token == "" {
		return "", "", "", errors.New("token is empty") // Возвращаем ошибку, если токен пустой
	}
	claims := &Claims{} // Создаем новую структуру для хранения данных из токена
	// Парсим токен с указанными claims и секретным ключом
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil // Возвращаем секретный ключ для проверки подписи
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", "", errors.New("signature is invalid") // Возвращаем ошибку, если подпись недействительна
		}
		return "", "", "", errors.New("token is invalid") // Возвращаем ошибку, если токен недействителен
	}
	if !parsedToken.Valid {
		return "", "", "", errors.New("parsed token is invalid") // Возвращаем ошибку, если разобранный токен недействителен
	}
	if claims == nil {
		return "", "", "", errors.New("token claims are nil") // Возвращаем ошибку, если поля токена не заполнены
	}
	// Возвращаем идентификатор, email и тип пользователя из токена
	return claims.Id, claims.Email, claims.UserType, nil
}
