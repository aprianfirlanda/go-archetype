package httpctx

import "github.com/gofiber/fiber/v2"

func GetUser(c *fiber.Ctx) map[string]any {
	if v := c.Locals("user"); v != nil {
		if user, ok := v.(map[string]any); ok {
			return user
		}
	}
	return nil
}

/*
HOW TO USE IN HANDLERS

user := httpctx.GetUser(c)

if user == nil {
	return fiber.NewError(fiber.StatusUnauthorized, "user not found")
}

username := user["preferred_username"]
roles := user["roles"]
permissions := user["permissions"]
resources := user["resources"]
*/
