package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.Lshortfile|log.LstdFlags)
	infoLog  = log.New(os.Stdout, "\033[33m[info ]\033[0m", log.Lshortfile|log.LstdFlags)
	loggers  = []*log.Logger{errorLog, infoLog}

	mu = sync.Mutex{}
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

func setLevel(level int) {

	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}
}
