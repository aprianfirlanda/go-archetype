package logging

import "github.com/sirupsen/logrus"

// ComponentLogger returns the logger with a fixed "component" field.
//
// Example component values:
//   - "http.server"
//   - "http.middleware.auth_jwt"
//   - "domain.user_service"
//   - "infra.db"
func ComponentLogger(logger *logrus.Entry, component string) *logrus.Entry {
	if logger == nil {
		logger = logrus.NewEntry(logrus.StandardLogger())
	}
	return logger.WithField("component", component)
}
