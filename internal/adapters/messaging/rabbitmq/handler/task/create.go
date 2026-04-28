package taskhandlermq

import (
	"context"
	"encoding/json"
	taskcmd "go-archetype/internal/application/task/command"
)

func (h *Handler) Create(ctx context.Context, payload []byte) error {
	var cmd taskcmd.Create
	if err := json.Unmarshal(payload, &cmd); err != nil {
		return err
	}

	_, err := h.service.Create(ctx, cmd)
	return err
}
