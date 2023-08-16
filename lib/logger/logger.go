package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LogCfg struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	timeFormat string `yaml:"time-format"`
}
type logLevel int

const (
	DEBUG logLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	F                  *os.File
	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func Setup(cfg *LogCfg) {
	var err error
	dir := cfg.Path
	filename := fmt.Sprintf("%s-%s.%s", cfg.Name, time.Now().Format(cfg.timeFormat), cfg.Ext)
	logFile, err := mustOpen(filename, dir)
	if err != nil {
		fmt.Errorf("logging.Setup err: %s", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(mw, DefaultPrefix, log.LstdFlags)
}

func setPrefix(level logLevel) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%s]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}

func Debug(v ...any) {
	setPrefix(DEBUG)
	logger.Println(v)
}

func Info(v ...any) {
	setPrefix(INFO)
	logger.Println(v)
}

func Warn(v ...any) {
	setPrefix(WARNING)
	logger.Println(v)
}

func Error(v ...any) {
	setPrefix(ERROR)
	logger.Println(v)
}

func Fatal(v ...any) {
	setPrefix(FATAL)
	logger.Fatalln(v)
}
