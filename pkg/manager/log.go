package manager

import "Kontroller/logging"

var log *logging.Logging

func init() {
	log = logging.NewLogging("manager")
}
