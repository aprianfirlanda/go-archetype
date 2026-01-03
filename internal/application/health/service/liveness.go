package healthsvc

import "context"

// Liveness is the app process alive?
func (s *service) Liveness(_ context.Context) bool {
	return true
}
