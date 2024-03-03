package backend

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// RequireAuthMiddleware - Middleware для проверки токена авторизации.
// Если токен отсутствует или неверен, возвращает ошибку аутентификации.
func RequireAuthMiddleware(c *gin.Context) {
	var token string
	authHeader := c.GetHeader("Authorization") // Получаем заголовок Authorization
	fields := strings.Fields(authHeader)       // Разбиваем на части по пробелу

	// Если заголовок не пустой и начинается с "Bearer", извлекаем токен
	if len(fields) != 0 && fields[0] == "Bearer" {
		token = fields[1]
	}

	// Если заголовок пустой или отсутствует токен, возвращаем ошибку аутентификации
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Попробуйте сначала войти",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Если токен недействителен, возвращаем ошибку аутентификации
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Проверяем токен на валидность и извлекаем данные пользователя
	id, email, role, err := VerifyToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Проверка токена не удалась"})
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Устанавливаем значения идентификатора, электронной почты и роли в контекст запроса
	c.Set("id", id)
	c.Set("email", email)
	c.Set("role", role)
	c.Next()
}

// isAdmin - Middleware для проверки роли пользователя.
// Если пользователь не является администратором, возвращает сообщение об ошибке.
func isAdmin(c *gin.Context) {
	role, exists := c.Get("role") // Получаем значение роли из контекста запроса
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Роль не найдена в контексте"})
		return
	}

	roleStr, ok := role.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось преобразовать роль в строку"})
		return
	}

	// Если роль пользователя - администратор, передаем управление следующему обработчику
	if roleStr == "ADMIN" {
		c.Next()
	} else {
		// В противном случае возвращаем сообщение о том, что пользователь не является администратором
		c.JSON(http.StatusOK, gin.H{"message": "Действия обычного пользователя"})
	}
}
