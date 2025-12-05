package logging

import "github.com/sirupsen/logrus"

// WithComponent returns a logger entry with a fixed "component" field.
//
// Example component values:
//   - "http.server"
//   - "http.middleware.auth_jwt"
//   - "domain.user_service"
//   - "infra.db"
func WithComponent(logger *logrus.Entry, component string) *logrus.Entry {
	return logger.WithField("component", component)
}

// WithComponentAndFields returns a logger entry with "component" and extra fields.
func WithComponentAndFields(logger *logrus.Entry, component string, fields logrus.Fields) *logrus.Entry {
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["component"] = component
	return logger.WithFields(fields)
}
