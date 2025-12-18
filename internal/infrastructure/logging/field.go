package logging

import "github.com/sirupsen/logrus"

func Field(key string, value any) logrus.Fields {
	return logrus.Fields{key: value}
}

func Fields(kv map[string]any) logrus.Fields {
	return logrus.Fields(kv)
}
