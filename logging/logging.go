package logging

import (
	"Kontroller/config"
	"flag"
	"fmt"
	"k8s.io/klog/v2"
	"strconv"
)

// Log level constants
const (
	LevelFatal int32 = iota + 1
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
)

// Map of log levels to klog.Level
var loggingLevel = map[int32]klog.Level{
	LevelFatal:   klog.Level(1),
	LevelError:   klog.Level(2),
	LevelWarning: klog.Level(3),
	LevelInfo:    klog.Level(4),
	LevelDebug:   klog.Level(5),
}
var level klog.Level

// Initialize klog
func init() {
	klog.InitFlags(nil)
	// Get log level
	level = loggingLevel[config.Cfg.Log.Level]
	// Ensure log level is within range
	if level < loggingLevel[LevelFatal] {
		level = loggingLevel[LevelFatal]
	}
	if level > loggingLevel[LevelDebug] {
		level = loggingLevel[LevelDebug]
	}
	// Set klog log level
	flag.Set("v", strconv.Itoa(int(level)))
}

// Logging struct
type Logging struct {
	prefix string     // Log prefix
	level  klog.Level // Log level
}

// Debugf method, output DEBUG level log
func (l *Logging) Debugf(format string, args ...interface{}) {
	klog.V(loggingLevel[LevelDebug]).Infof("[DEBUG][%s]%s", l.prefix, fmt.Sprintf(format, args...))
}

// Infof method, output INFO level log
func (l *Logging) Infof(format string, args ...interface{}) {
	klog.V(loggingLevel[LevelInfo]).Infof("[INFO][%s]%s", l.prefix, fmt.Sprintf(format, args...))
}

// Warnf method, output WARN level log
func (l *Logging) Warnf(format string, args ...interface{}) {
	klog.V(loggingLevel[LevelWarning]).Infof("[WARN][%s]%s", l.prefix, fmt.Sprintf(format, args...))
}

// Errorf method, output ERROR level log
func (l *Logging) Errorf(format string, args ...interface{}) {
	klog.V(loggingLevel[LevelError]).Infof("[ERROR][%s]%s", l.prefix, fmt.Sprintf(format, args...))
}

// Fatalf method, output FATAL level log
func (l *Logging) Fatalf(format string, args ...interface{}) {
	klog.V(loggingLevel[LevelFatal]).Infof("[FATAL][%s]%s", l.prefix, fmt.Sprintf(format, args...))
}

// NewLogging function, create Logging instance
func NewLogging(prefix string) *Logging {
	return &Logging{prefix: prefix, level: level}
}
