package watcher

import (
	"fmt"
	"os"
)

func addStatusForHealthCheck(msg []byte) error {
	healthCheckFile, openErr := os.OpenFile(HealthCheckFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return fmt.Errorf("error opening %s: %w", HealthCheckFile, openErr)
	}
	if _, writeErr := healthCheckFile.Write(msg); writeErr != nil {
		closeErr := healthCheckFile.Close()
		if closeErr != nil {
			errorMsg := "error closing %s: %w -  after error writing: %s"
			return fmt.Errorf(errorMsg, HealthCheckFile, closeErr, writeErr.Error())
		}
		return fmt.Errorf("error writing watcher startup to %s: %w", HealthCheckFile, writeErr)
	}
	if closeErr := healthCheckFile.Close(); closeErr != nil {
		return fmt.Errorf("error closing %s: %w", HealthCheckFile, closeErr)
	}
	return nil
}
