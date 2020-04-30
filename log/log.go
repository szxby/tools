// log

package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	LDebug = iota
	LInfo
	LWarn
	LError
	LFatal
	maxLogCycle = 15
)

var (
	enableDebug              = true
	defaultLogPath           = getDefaultLogPath()
	defaultLogLevel          = LDebug
	fileLog         *FileLog = NewFileLog(defaultLogPath, defaultLogLevel)
	LogTag                   = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

// default path
func getDefaultLogPath() string {
	procPath := string(os.Args[0])
	n := strings.LastIndexByte(procPath, os.PathSeparator)
	procName := procPath[n+1:]
	return fmt.Sprintf("log%c%s%crun.log", os.PathSeparator, procName, os.PathSeparator)
}

type FileLog struct {
	path     string
	MaxCycle int
	Level    int
	t        time.Time
	f        *os.File
	logger   *log.Logger
	mu       sync.Mutex
}

func NewFileLog(path string, level int) *FileLog {
	l := &FileLog{
		path:  path,
		Level: level,
	}
	l.t, _ = time.Parse("2006-01-01", "1900-01-01")
	f, _ := os.Open(os.DevNull)
	l.logger = log.New(f, "", log.Lshortfile|log.LstdFlags)
	if enableDebug {
		log.SetOutput(os.Stdout)
		log.SetFlags(log.Lshortfile | log.LstdFlags)
	}
	return l
}

func (l *FileLog) NewFile(new_path string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f != nil {
		l.f.Close()
		os.Rename(l.path, new_path)
	}
	now := time.Now()
	n := strings.LastIndexByte(l.path, os.PathSeparator)
	if n >= 0 {
		dir := l.path[:n]
		os.MkdirAll(dir, 0755)
	}
	deadline := now.Add(-time.Duration(maxLogCycle) * time.Hour * 24)
	deadlinePath := fmt.Sprintf("%s.%02d-%02d", l.path, deadline.Month(), deadline.Day())
	// fmt.Println(deadlinePath)
	if _, err := os.Lstat(deadlinePath); err == nil {
		os.Remove(deadlinePath)
	}
	f, err := os.OpenFile(l.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		panic("create log file error")
		return err
	}
	l.f = f
	l.t = now
	l.logger.SetOutput(f)
	return nil
}

func (l *FileLog) Output(level int, s string) {
	if level < LDebug || level > LFatal {
		return
	}
	if level < l.Level {
		return
	}
	now := time.Now()
	if l.t.YearDay() != now.YearDay() {
		newpath := fmt.Sprintf("%s.%02d-%02d", l.path, l.t.Month(), l.t.Day())
		l.NewFile(newpath)
		l.t = now
	}
	if l.logger == nil {
		panic("Log output path unset")
	}
	s = fmt.Sprintf("[%s] %s", LogTag[level], s)
	l.logger.Output(3, s)
	if enableDebug {
		log.Output(3, s)
	}
}

func SetLevel(level int) {
	if level < LDebug || level > LFatal {
		return
	}
	fileLog.Level = level
}

func SetOutput(path string) {
	if fileLog.path == path {
		return
	}
	fileLog.NewFile(path)
	fileLog.path = path
}

func Debug(format string, v ...interface{}) {
	fileLog.Output(LDebug, fmt.Sprintf(format, v...))
}

// func Debug(v ...interface{}) {
// 	fileLog.Output(LDebug, fmt.Sprintln(v...))
// }

func Release(format string, v ...interface{}) {
	fileLog.Output(LInfo, fmt.Sprintf(format, v...))
}

// func Info(v ...interface{}) {
// 	fileLog.Output(LInfo, fmt.Sprintln(v...))
// }

func Warn(format string, v ...interface{}) {
	fileLog.Output(LWarn, fmt.Sprintf(format, v...))
}

// func Warn(v ...interface{}) {
// 	fileLog.Output(LWarn, fmt.Sprintln(v...))
// }

func Error(format string, v ...interface{}) {
	fileLog.Output(LError, fmt.Sprintf(format, v...))
}

// func Error(v ...interface{}) {
// 	fileLog.Output(LError, fmt.Sprintln(v...))
// }

func Fatal(format string, v ...interface{}) {
	fileLog.Output(LFatal, fmt.Sprintf(format, v...))
	os.Exit(0)
}

// func Fatal(v ...interface{}) {
// 	fileLog.Output(LFatal, fmt.Sprintln(v...))
// 	os.Exit(0)
// }
