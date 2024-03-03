package backend

import "github.com/sirupsen/logrus"

// log представляет собой экземпляр логгера Logrus.
var log = logrus.New()

// GetLogger возвращает экземпляр логгера.
func GetLogger() *logrus.Logger {
	return log
}
