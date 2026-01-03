package healthsvc

import "context"

// Readiness is the app ready to receive traffic?
func (s *service) Readiness(ctx context.Context) error {
	if s.db != nil {
		if err := s.db.Ping(ctx); err != nil {
			return err
		}
	}

	// future checks:
	// cache.Ping
	// mq.Ping

	return nil
}
