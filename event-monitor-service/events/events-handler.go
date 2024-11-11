package events

import (
	"github.com/labstack/gommon/log"
)

// TriggerOnboarding dummy function
func TriggerOnboarding() {
	log.Info("Onboarding has been triggered")
}

// RefreshConfiguration dummy function
func RefreshConfiguration() {
	log.Info("Configuration has been refreshed")
}

// UpdateSubsystems dummy function
func UpdateSubsystems() {
	log.Info("Subsystems have been updated")
}

// Shutdown dummy function
func Shutdown() {
	log.Info("Shutdown has been triggered")
}

// Reinstall dummy function
func Reinstall() {
	log.Info("Reinstall has been triggered")
}
