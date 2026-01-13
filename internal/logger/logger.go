package logger

import (
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
	error *log.Logger
	file *os.File
}

func New() (*Logger, error) {
	// create logs directory if not exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		"logs/app.text",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)

	if err != nil {
		return nil, err
	}

	return &Logger{
		info: log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		file: file,
	}, nil
}

func (l *Logger) Info(msg string) {
	l.info.Println(msg)
}

func (l *Logger) Error(msg string, err any) {
	l.error.Printf("$s: %v\n", msg, err)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
