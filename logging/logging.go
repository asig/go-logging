package logging

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Level int8

const (
	FATAL Level = iota
	SEVERE
	WARNING
	INFO
	DEBUG
)

type Logger struct {
	name       string
	underlying *log.Logger
}

var (
	stringToLevel = map[string]Level{
		"FATAL":   FATAL,
		"SEVERE":  SEVERE,
		"WARNING": WARNING,
		"INFO":    INFO,
		"DEBUG":   DEBUG,
	}
	levelToString = map[Level]string{
		FATAL:   "FATAL",
		SEVERE:  "SEVERE",
		WARNING: "WARNING",
		INFO:    "INFO",
		DEBUG:   "DEBUG",
	}

	levelPerComponent map[string]Level

	loggerMap     = make(map[string]*Logger)
	loggerMapLock sync.Mutex

	flagLogFile  = flag.String("logfile", "", "filename of the log file. if empty, log will be written to stderr")
	flagLogLevel = flag.String("loglevel", "FATAL", "minimal level that is logged")

	initialized bool = false
	logFile     *os.File
	logLevel    Level
)

func Initialize() {
	parseLogLevelSpec(*flagLogLevel)

	logFile = os.Stderr
	if len(*flagLogFile) > 0 {
		logFile, _ = os.Create(*flagLogFile)
	}
	initialized = true
}

func parseLogLevelSpec(levelSpec string) {
	levelPerComponent = map[string]Level{
		"": WARNING,
	}
	for _, spec := range strings.Split(levelSpec, ",") {
		component, rawLevel := "", ""
		parts := strings.Split(spec, "=")
		if len(parts) == 1 {
			rawLevel = parts[0]
		} else {
			component = strings.TrimSpace(parts[0])
			rawLevel = parts[1]
		}
		level, levelFound := stringToLevel[strings.TrimSpace(rawLevel)]
		if !levelFound {
			panic("Unknown log level " + rawLevel)
		}
		levelPerComponent[component] = level
	}
}

func Get(name string) *Logger {
	if !initialized {
		panic("Call Initialize() first")
	}
	loggerMapLock.Lock()
	defer loggerMapLock.Unlock()
	logger, found := loggerMap[name]
	if !found {
		logger = &Logger{
			name,
			log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds)}
		loggerMap[name] = logger
	}
	return logger
}

func (self *Logger) Debug(v ...interface{}) {
	self.Log(DEBUG, v...)
}

func (self *Logger) Debugf(format string, v ...interface{}) {
	self.Logf(DEBUG, format, v...)
}

func (self *Logger) Info(v ...interface{}) {
	self.Log(INFO, v...)
}

func (self *Logger) Infof(format string, v ...interface{}) {
	self.Logf(INFO, format, v...)
}

func (self *Logger) Warning(v ...interface{}) {
	self.Log(WARNING, v...)
}

func (self *Logger) Warningf(format string, v ...interface{}) {
	self.Logf(WARNING, format, v...)
}

func (self *Logger) Severe(v ...interface{}) {
	self.Log(SEVERE, v...)
}

func (self *Logger) Severef(format string, v ...interface{}) {
	self.Logf(SEVERE, format, v...)
}

func (self *Logger) Fatal(v ...interface{}) {
	self.Log(FATAL, v...)
	os.Exit(1)
}

func (self *Logger) Fatalf(format string, v ...interface{}) {
	self.Logf(FATAL, format, v...)
	os.Exit(1)
}

func (self *Logger) Logf(level Level, format string, v ...interface{}) {
	if level <= getComponentLevel(self.name) {
		message := fmt.Sprintf(format, v...)
		self.underlying.Printf("%s %-8s %s", levelToString[level], self.name, message)
	}
}

func (self *Logger) Log(level Level, v ...interface{}) {
	if level <= getComponentLevel(self.name) {
		message := fmt.Sprint(v...)
		self.underlying.Printf("%s %-8s %s", levelToString[level], self.name, message)
	}
}

func getComponentLevel(component string) Level {
	level, levelFound := levelPerComponent[component]
	if levelFound {
		return level
	}
	pos := strings.LastIndex(component, "/")
	if pos > 0 {
		return getComponentLevel(component[:pos])
	} else {
		// no component to strip... return default value
		return getComponentLevel("")
	}
}
