package jwtutil

func ExtractPermissions(user map[string]interface{}) ([]string, error) {
	permissionsRaw, ok := user["permissions"].([]interface{})
	if !ok || len(permissionsRaw) == 0 {
		return nil, ErrInvalidPermissions
	}

	permissions := make([]string, 0, len(permissionsRaw))
	for _, r := range permissionsRaw {
		role, ok := r.(string)
		if !ok {
			continue
		}
		permissions = append(permissions, role)
	}
	return permissions, nil
}
