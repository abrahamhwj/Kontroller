package logging

import (
	"testing"
)

func TestLogging_Debugf(t *testing.T) {
	log := NewLogging("test")
	log.Debugf("debug message")
}
func TestLogging_Infof(t *testing.T) {
	log := NewLogging("test")
	log.Infof("info message")
}
func TestLogging_Warnf(t *testing.T) {
	log := NewLogging("test")
	log.Warnf("warn message")
}
func TestLogging_Errorf(t *testing.T) {
	log := NewLogging("test")
	log.Errorf("error message")
}
func TestLogging_Fatalf(t *testing.T) {
	log := NewLogging("test")
	log.Fatalf("fatal message")
}
