package jwtutil

func HasPermission(permissions []string, target string) bool {
	for _, r := range permissions {
		if r == target {
			return true
		}
	}
	return false
}
