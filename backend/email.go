package backend

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/smtp"
)

// GenerateVerificationCode генерирует случайный код верификации из шести цифр.
// Код формируется путем выбора случайных цифр из набора "0123456789".
// Возвращает строку с сгенерированным кодом верификации.
func GenerateVerificationCode() string {
	// Генерация случайного кода из шести цифр
	code := make([]byte, 6)
	characters := "0123456789"
	for i := 0; i < 6; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		code[i] = characters[randomIndex.Int64()]
	}
	return string(code)
}

// SendVerificationCodeEmail отправляет электронное письмо с кодом верификации на указанный адрес электронной почты.
// Параметры:
// - email: адрес электронной почты, на который отправляется письмо.
// - verificationCode: сгенерированный код верификации.
// Возвращает ошибку, если отправка письма не удалась.
func SendVerificationCodeEmail(email, verificationCode string) error {
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: Verification Code\n\nYour verification code is: %s", "tilesjan2005@gmail.com", email, verificationCode)
	auth := smtp.PlainAuth("", "tilesjan2005@gmail.com", "hwlq wpyq huhd zcuy", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "tilesjan2005@gmail.com", []string{email}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

// SendSpamForUser отправляет спам-сообщение пользователю с указанным адресом электронной почты.
// Параметры:
// - email: адрес электронной почты пользователя.
// - text: текст сообщения, который будет отправлен пользователю.
// Возвращает ошибку, если отправка письма не удалась.
func SendSpamForUser(email, text string) error {
	// Отправка спам-сообщения пользователю через SMTP
	msg := []byte(text)
	auth := smtp.PlainAuth("", "tilesjan2005@gmail.com", "hwlq wpyq huhd zcuy", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "tilesjan2005@gmail.com", []string{email}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}

func SendSpamForAdmin(email string, book Book) error { //idk
	bookJSON, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("failed to marshal book to JSON: %v", err)
	}

	// Отправка спам-сообщения пользователю через SMTP
	auth := smtp.PlainAuth("", "your-email@gmail.com", "your-password", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "your-email@gmail.com", []string{email}, bookJSON)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
