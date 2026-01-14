package logger

import (
	"io"
	"log"
	"os"
)

type Logger struct {
	info  *log.Logger
	error *log.Logger
	file  *os.File
}

func New() (*Logger, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		"logs/app.txt",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return nil, err
	}
	infoWriter := io.MultiWriter(os.Stdout, file)
	errWriter := io.MultiWriter(os.Stderr, file)

	return &Logger{
		info:  log.New(infoWriter, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile),
		error: log.New(errWriter, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile),
		file:  file,
	}, nil
}

func (l *Logger) Info(msg string) {
	l.info.Println(msg)
}

func (l *Logger) Error(msg string, err any) {
	l.error.Printf("%s: %v\n", msg, err)
}

func (l *Logger) Close() error {
	return l.file.Close()
}
