package handlers

import "github.com/gofiber/fiber/v2"

func (r *Resolver) send(c *fiber.Ctx) error {
	var pushPresenter PushPresenter
	err := c.BodyParser(&pushPresenter)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			ErrorPresenter{
				Reason: "error with parsing body: " + err.Error(),
			},
		)
	}

	push := pushPresenter.ToCore()

	pushID, err := r.pushService.SendPush(c.UserContext(), push)
	if err != nil {
		resp, status := mapError(err)

		return c.Status(status).JSON(resp)
	}

	return c.Status(fiber.StatusOK).JSON(
		SendResponsePresenter{
			PushID: pushID.String(),
		},
	)
}
