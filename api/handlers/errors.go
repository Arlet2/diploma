package handlers

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func mapError(err error) (resp ErrorPresenter, status int) {
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		var reason strings.Builder

		for _, err := range validationErrs {
			reason.WriteString(fmt.Sprintf("validation error: %s. ", err.Error()))
		}

		return ErrorPresenter{
			Reason: reason.String(),
		}, fiber.StatusUnprocessableEntity
	}

	reason := err.Error()

	return ErrorPresenter{
		Reason: reason,
	}, fiber.StatusInternalServerError
}
