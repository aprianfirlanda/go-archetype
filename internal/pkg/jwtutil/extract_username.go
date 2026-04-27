package jwtutil

func ExtractUsername(user map[string]interface{}) (string, error) {
	username, ok := user["preferred_username"].(string)
	if !ok || username == "" {
		return "", ErrInvalidUsername
	}
	return username, nil
}
