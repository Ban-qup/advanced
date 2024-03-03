package backend

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Email    string `gorm:"unique"`
	Photo    string
	Password string
	Role     string
	AboutMe  string
}

type Chapter struct {
	gorm.Model
	Name  string
	Words string
}

type Book struct {
	gorm.Model
	Name             string `gorm:"unique"`
	Photo            string
	BriefInformation string
	Genre            string
	Author           string
	Translator       *User `gorm:"belongsTo:book_translator"`
	Finished         bool  `gorm:"notnull"`
	ChapterQuantity  uint
	Chapters         []*Chapter `gorm:"hasMany:book_chapters"`
}

type VerificationCode struct {
	ID             uint      `gorm:"primaryKey"`
	Email          string    `gorm:"uniqueIndex"`
	Code           string    `gorm:"not null"`
	ExpirationTime time.Time `gorm:"not null"`
}
