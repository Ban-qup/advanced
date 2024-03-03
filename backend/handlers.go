package backend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func register(c *gin.Context) {
	// Начало процесса регистрации пользователя
	GetLogger().Info("Starting user registration")

	// Попытка привязать JSON-данные запроса к структуре User
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		// Обработка ошибки невалидного запроса
		GetLogger().Error("Invalid registration request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	// Проверка, существует ли уже аккаунт с указанным email
	if err := db.Where("email = ?", newUser.Email).First(&newUser).Error; err == nil {
		// Обработка случая, когда аккаунт уже существует
		fmt.Println(err)
		GetLogger().Error("Account already registered for email:", newUser.Email)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "The account is already registered"})
		return
	}

	// Определение роли пользователя на основе email и пароля
	if newUser.Password == "123" && newUser.Email == "tilesjan2005@gmail.com" {
		newUser.Role = "ADMIN"
	} else {
		newUser.Role = "USER"
	}

	// Хэширование пароля перед сохранением в базе данных
	fmt.Println(newUser.Password)
	hashedPassword, _ := HashPassword(newUser.Password)
	newUser.Password = hashedPassword
	fmt.Println(newUser.Password)

	// Создание записи пользователя в базе данных
	db.Create(&newUser)

	// Создание JWT-токена для пользователя
	signedToken, _ := CreateToken(strconv.Itoa(int(newUser.ID)), newUser.Email, newUser.Role)

	// Установка куки с JWT-токеном для пользователя
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	// Логирование успешной регистрации пользователя
	GetLogger().Info("User registered successfully")

	// Отправка ответа об успешной регистрации
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
func login(c *gin.Context) {
	// Начало процесса входа пользователя в систему
	GetLogger().Info("Starting user login")

	// Создание структуры для хранения данных формы входа
	var user User
	type loginForm struct {
		Username string `json:"username"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password"`
	}

	// Привязка JSON-данных запроса к структуре формы входа
	var newUser loginForm
	if err := c.ShouldBindJSON(&newUser); err != nil {
		// Обработка ошибки невалидного запроса
		GetLogger().Error("Invalid registration request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(newUser)

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	// Поиск пользователя по email в базе данных
	if err := db.Where("email = ?", newUser.Email).First(&user).Error; err != nil {
		// Обработка ошибки неверного email или пароля
		GetLogger().Error("Invalid email or password", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	fmt.Println(user)
	fmt.Println(HashPassword(newUser.Password))
	fmt.Println((user.Password))

	// Проверка соответствия хэшированного пароля
	if !CheckPasswordHash(newUser.Password, user.Password) {
		// Обработка ошибки аутентификации
		GetLogger().Error("Authentication failed for user:", user.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Определение роли пользователя на основе email и пароля
	if (user.Password == "Ali12" && user.Email == "aa5331865@gmail.com") || (user.Password == "Tim12" && user.Email == "timchik.mux@mail.ru") {
		user.Role = "ADMIN"
	} else {
		user.Role = "USER"
	}
	token, err := CreateToken(strconv.Itoa(int(user.ID)), user.Email, user.Role)
	if err != nil {
		// Обработка ошибки создания токена
		GetLogger().Error("Failed to create token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token", "data": token})
		return
	}

	// Установка куки с JWT-токеном для пользователя
	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, &cookie)

	// Логирование успешного входа пользователя
	GetLogger().Info("User login successful")

	// Отправка ответа об успешном входе и передача токена
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Getvcode(context *gin.Context) {
	// Начало процесса восстановления пароля
	GetLogger().Info("Starting forgot password process")

	// Структура для хранения данных формы запроса на восстановление пароля
	type form struct {
		Email string `json:"email"`
	}

	// Привязка JSON-данных запроса к структуре формы
	var forminput form
	if err := context.BindJSON(&forminput); err != nil {
		// Обработка невалидного запроса на восстановление пароля
		GetLogger().Error("Invalid forgot password request:", err)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	fmt.Println(forminput)

	// Генерация кода подтверждения
	verificationCode := GenerateVerificationCode()

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		GetLogger().Error("Failed to connect to the database:", err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	// Сохранение кода подтверждения в базе данных
	err = SaveVerificationCode(context, db, forminput.Email, verificationCode)
	if err != nil {
		// Обработка ошибки сохранения кода подтверждения в базе данных
		GetLogger().Error("Failed to save verification code to Redis:", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправка кода подтверждения на электронную почту
	err = SendVerificationCodeEmail(forminput.Email, verificationCode)

	// Логирование успешной отправки кода подтверждения на электронную почту
	GetLogger().Info("Verification code sent to email successfully")

	// Отправка ответа об успешной отправке кода подтверждения на электронную почту
	context.JSON(http.StatusOK, gin.H{"message": "Verification code sent to your email"})
}
func Checkvcode(c *gin.Context) {
	// Проверка кода подтверждения
	GetLogger().Info("Checking verification code")

	// Структура для проверки кода подтверждения
	type CheckCode struct {
		Email string `json:"email,omitempty"`
		Code  string `json:"code"`
	}

	// Привязка JSON-данных запроса к структуре для проверки кода
	var code CheckCode
	if err := c.BindJSON(&code); err != nil {
		// Обработка невалидного запроса на проверку кода
		GetLogger().Error("Invalid code check request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		GetLogger().Error("Failed to connect to the database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database connection failed"})
		return
	}

	// Проверка кода подтверждения
	validCode, err := CheckVerificationCode(c, db, code.Email, code.Code)
	if err != nil {
		// Обработка ошибки проверки кода подтверждения
		GetLogger().Error("Failed to verify code:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify code"})
		return
	}

	// Если код не действителен, возврат ошибки
	if !validCode {
		GetLogger().Error("Failed to save user information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user information"})
		return
	}

	// Успешная проверка кода
	GetLogger().Info("successful")
	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
func UpdateUser(c *gin.Context) {
	// Обновление данных пользователя
	GetLogger().Info("Starting user update")

	// Получение электронной почты пользователя из контекста
	userEmail, emailExists := c.Get("email")
	if !emailExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	// Преобразование электронной почты в строку
	email, ok := userEmail.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert email to string"})
		return
	}

	// Получение данных пользователя из запроса
	var updateUser User
	if err := c.ShouldBindJSON(&updateUser); err != nil {
		GetLogger().Error("Invalid user update request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		GetLogger().Error("Failed to connect to the database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Поиск существующего пользователя в базе данных
	var existingUser User
	if err := db.Where("email = ?", email).First(&existingUser).Error; err != nil {
		// Обработка ошибки поиска пользователя
		GetLogger().Error("Failed to find user:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Обновление полей пользователя, если они не пустые
	if updateUser.Username != "" {
		existingUser.Username = updateUser.Username
	}
	if updateUser.Email != "" {
		existingUser.Email = updateUser.Email
	}
	if updateUser.Password != "" {
		// Обновление пароля (возможно, стоит хешировать перед сохранением)
		existingUser.Password = updateUser.Password
	}
	if updateUser.AboutMe != "" {
		existingUser.AboutMe = updateUser.AboutMe
	}
	if updateUser.Photo != "" {
		existingUser.Photo = updateUser.Photo
	}
	// Обновление других полей по необходимости

	// Сохранение обновленной информации о пользователе в базе данных
	if err := db.Save(&existingUser).Error; err != nil {
		// Обработка ошибки обновления пользователя
		GetLogger().Error("Failed to update user:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Успешное обновление данных пользователя
	GetLogger().Info("User updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": existingUser})
}
func seeProfile(c *gin.Context) {
	// Получение электронной почты пользователя, информацию о котором нужно просмотреть, из параметров запроса
	userEmail := c.Query("email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User email not provided"})
		return
	}

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		GetLogger().Error("Failed to connect to the database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Поиск пользователя в базе данных по электронной почте
	var user User
	if err := db.Where("email = ?", userEmail).First(&user).Error; err != nil {
		// Обработка ошибки поиска пользователя
		GetLogger().Error("Failed to find user:", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Создание копии пользователя без поля пароля
	profile := User{
		Username: user.Username,
		Email:    user.Email,
		Photo:    user.Photo,
		Role:     user.Role,
		AboutMe:  user.AboutMe,
	}

	// Отправка информации о пользователе в качестве ответа
	c.JSON(http.StatusOK, gin.H{"user": profile})
}
func createBook(c *gin.Context) {
	// Получение email пользователя из запроса
	userEmail := c.Query("email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User email not provided"})
		return
	}

	// Получение данных о текущем пользователе из контекста запроса
	userRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
		return
	}

	// Проверка роли пользователя
	roleStr, ok := userRole.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert role to string"})
		return
	}

	if roleStr != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can create books"})
		return
	}

	// Попытка привязать JSON-данные запроса к структуре Book
	var newBook Book
	if err := c.ShouldBindJSON(&newBook); err != nil {
		// Обработка ошибки невалидного запроса
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка заполненности всех полей при создании книги
	if newBook.Name == "" || newBook.Photo == "" || newBook.BriefInformation == "" || newBook.Genre == "" || newBook.Author == "" || newBook.Translator == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields must be provided when creating a book"})
		return
	}

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Создание записи книги в базе данных
	if err := db.Create(&newBook).Error; err != nil {
		// Обработка ошибки создания книги
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	// Отправка ответа об успешном создании книги
	c.JSON(http.StatusOK, gin.H{"message": "Book created successfully", "book": newBook})
}
func addChapter(c *gin.Context) {
	// Получение данных о текущем пользователе из контекста запроса
	userEmail, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	// Получение данных о книге из запроса
	bookID := c.Param("id")

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		// Обработка ошибки подключения к базе данных
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Поиск книги в базе данных
	var book Book
	if err := db.First(&book, bookID).Error; err != nil {
		// Обработка ошибки поиска книги
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Проверка, является ли текущий пользователь переводчиком книги
	if book.Translator == nil || book.Translator.Email != userEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only the translator of the book can add chapters"})
		return
	}

	// Попытка привязать JSON-данные запроса к структуре Chapter
	var newChapter Chapter
	if err := c.ShouldBindJSON(&newChapter); err != nil {
		// Обработка ошибки невалидного запроса
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка обязательных полей главы
	if newChapter.Name == "" || newChapter.Words == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chapter name and words must be provided"})
		return
	}

	// Добавление главы к книге
	db.Model(&book).Association("Chapters").Append(&newChapter)

	// Отправка ответа об успешном добавлении главы
	c.JSON(http.StatusOK, gin.H{"message": "Chapter added successfully", "chapter": newChapter})
}
func UpdateBook(c *gin.Context) {
	// Получение электронной почты пользователя из контекста
	_, emailExists := c.Get("email")
	if !emailExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User email not found"})
		return
	}

	// Проверка роли пользователя (только администраторы могут обновлять книги)
	userRole, roleExists := c.Get("role")
	if !roleExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User role not found"})
		return
	}

	// Проверка роли пользователя
	roleStr, ok := userRole.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert role to string"})
		return
	}
	if roleStr != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only administrators can update books"})
		return
	}

	// Получение ID книги из параметров запроса
	bookID := c.Param("id")

	// Получение данных книги из запроса
	var updatedBook Book
	if err := c.ShouldBindJSON(&updatedBook); err != nil {
		GetLogger().Error("Invalid book update request:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		GetLogger().Error("Failed to connect to the database:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}

	// Проверка существования книги в базе данных
	var existingBook Book
	if err := db.First(&existingBook, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Обновление данных книги
	if updatedBook.Name != "" {
		existingBook.Name = updatedBook.Name
	}
	if updatedBook.Photo != "" {
		existingBook.Photo = updatedBook.Photo
	}
	if updatedBook.BriefInformation != "" {
		existingBook.BriefInformation = updatedBook.BriefInformation
	}
	if updatedBook.Genre != "" {
		existingBook.Genre = updatedBook.Genre
	}
	if updatedBook.Author != "" {
		existingBook.Author = updatedBook.Author
	}
	if updatedBook.Translator != nil {
		existingBook.Translator = updatedBook.Translator
	}
	if updatedBook.Finished {
		existingBook.Finished = updatedBook.Finished
	}
	// Сохранение обновленных данных в базе данных
	if err := db.Save(&existingBook).Error; err != nil {
		GetLogger().Error("Failed to update book:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	// Успешное обновление данных книги
	GetLogger().Info("Book updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully", "book": existingBook})
}
func deleteBook(c *gin.Context) {
	// Проверка электронной почты администратора
	adminEmail := "admin@example.com" // Замените на реальный адрес электронной почты администратора
	userEmail := c.Query("email")

	if userEmail != adminEmail {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to delete books"})
		return
	}

	// Проверка роли пользователя
	role, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Role not found in context"})
		return
	}

	// Проверка, является ли пользователь администратором
	if role.(string) != "ADMIN" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Only admins can delete books"})
		return
	}

	// Получение ID книги из параметров запроса
	bookID := c.Param("id")
	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	// Поиск книги в базе данных по ID
	var book Book
	if err := db.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Удаление книги из базы данных
	if err := db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	// Ответ об успешном удалении книги
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
func readBook(c *gin.Context) {
	// Получение ID книги из параметров запроса
	bookID := c.Param("id")

	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	// Находим книгу по её ID
	var book Book
	if err := db.First(&book, bookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}

	// Возвращаем информацию о книге
	c.JSON(http.StatusOK, gin.H{"book": book})
}
func readAllBooks(c *gin.Context) {
	// Подключение к базе данных
	db, err := dbConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	// Получение всех книг из базы данных
	var books []Book
	if err := db.Find(&books).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books from the database"})
		return
	}

	// Возвращаем список всех книг
	c.JSON(http.StatusOK, gin.H{"books": books})
}
func sendSpam(c *gin.Context) {
	// Получение всех пользователей из базы данных
	db, err := dbConnect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}

	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users from the database"})
		return
	}

	// Текст сообщения о новой книге
	message := "Уважаемые пользователи, рады сообщить вам о выходе новой книги!"

	// Отправка сообщения о новой книге каждому пользователю
	for _, user := range users {
		if err := SendSpamForUser(user.Email, message); err != nil {
			fmt.Printf("Failed to send spam email to %s: %v\n", user.Email, err)
			continue
		}
		fmt.Printf("Spam email sent successfully to %s\n", user.Email)
	}

	// Возвращаем сообщение об успешной отправке спама
	c.JSON(http.StatusOK, gin.H{"message": "Spam email sent successfully to all users"})
}
