package logger

import (
	"log"
	"os"
)

type Logger struct {
	Err  *log.Logger
	Info *log.Logger
}

func InitLogger() *Logger {

	logger := Logger{}

	logger.Info = log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	logger.Err = log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime)

	logger.Info.Print("Executing InitLogger.")
	return &logger
}
