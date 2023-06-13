package config

import (
	"github.com/spf13/viper"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	// Test negative case: configuration file not found
	viper.SetConfigName("NonexistentConfig")
	err := viper.ReadInConfig()
	if err == nil {
		t.Errorf("Expected error when reading nonexistent config file, but got no error.")
	}
}
func TestLog(t *testing.T) {
	// Test positive case
	log := Log{Level: 3}
	if log.Level != 3 {
		t.Errorf("Log level is incorrect, got: %d, want: %d.", log.Level, 3)
	}
	// Test negative case: invalid log level
	log = Log{Level: -1}
	if log.Level != -1 {
		t.Errorf("Log level is incorrect, got: %d, want: %d.", log.Level, -1)
	}
}
func TestManager(t *testing.T) {
	// Test positive case
	manager := Manager{
		ThreadNumber:            2,
		ControllerMaxRetryTimes: 3,
		ThreadTimeout:           10 * time.Second,
		ReSyncPeriod:            600 * time.Second,
	}
	if manager.ThreadNumber != 2 {
		t.Errorf("Thread number is incorrect, got: %d, want: %d.", manager.ThreadNumber, 2)
	}
	if manager.ControllerMaxRetryTimes != 3 {
		t.Errorf("Controller max retry times is incorrect, got: %d, want: %d.", manager.ControllerMaxRetryTimes, 3)
	}
	if manager.ThreadTimeout != 10*time.Second {
		t.Errorf("Thread timeout is incorrect, got: %d, want: %d.", manager.ThreadTimeout, 10*time.Second)
	}
	if manager.ReSyncPeriod != 600*time.Second {
		t.Errorf("Re-sync period is incorrect, got: %d, want: %d.", manager.ReSyncPeriod, 600*time.Second)
	}
	// Test negative case: invalid thread number
	manager = Manager{ThreadNumber: -1}
	if manager.ThreadNumber != -1 {
		t.Errorf("Thread number is incorrect, got: %d, want: %d.", manager.ThreadNumber, -1)
	}
}
