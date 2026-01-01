package demohandler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *DemoHandler) Panic(_ *fiber.Ctx) error {
	h.log.Warn("About to panic with nil pointer")

	var x *int
	_ = *x
	return nil
}
