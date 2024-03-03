package backend

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func rateLimiterMiddleware(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	// Выполнение автомиграции базы данных перед запуском сервера
	autoMigrate()

	// Создание лимитера скорости, разрешающего не более 3 запросов в секунду
	limiter := rate.NewLimiter(1, 3)

	r := gin.Default()

	// Применение middleware для ограничения скорости запросов ко всем маршрутам
	r.Use(rateLimiterMiddleware(limiter))

	// Установка обработчиков маршрутов для различных эндпоинтов
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/getvcode", Getvcode)
	r.POST("/checkvcode", Checkvcode)
	r.POST("/update", RequireAuthMiddleware, UpdateUser)
	r.GET("/profile", RequireAuthMiddleware, seeProfile)
	r.POST("/books/create", RequireAuthMiddleware, isAdmin, createBook)
	r.POST("/books/:id/chapters", RequireAuthMiddleware, addChapter)
	r.PUT("/books/:id/update", RequireAuthMiddleware, isAdmin, UpdateBook)
	r.DELETE("/books/:id/delete", RequireAuthMiddleware, isAdmin, deleteBook)
	r.GET("/books/:id/read", RequireAuthMiddleware, readBook)
	r.GET("/index", RequireAuthMiddleware, readAllBooks)
	r.POST("/sendSpam", RequireAuthMiddleware, isAdmin, sendSpam)

	// Создание HTTP-сервера с настройками адреса и обработчика маршрутов
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Запуск HTTP-сервера в горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("ListenAndServe: %v\n", err)
		}
	}()

	// Ожидание сигнала прерывания (например, Ctrl+C) для завершения работы сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutting down server...")

	// Определение контекста с таймаутом для graceful shutdown сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Выполнение graceful shutdown сервера с ожиданием завершения всех запросов
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
	fmt.Println("Server exiting")

}
