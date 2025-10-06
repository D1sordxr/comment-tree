package httputil

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type validationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func BindWithValidation(ctx *gin.Context, obj any) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return fmt.Errorf("binding failed: %w", err)
	}

	if err := validate.Struct(obj); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errs []validationError
			for _, fieldError := range validationErrors {
				errs = append(errs, validationError{
					Field:   fieldError.Field(),
					Message: fmt.Sprintf("Field %s failed validation: %s", fieldError.Field(), fieldError.Tag()),
				})
			}
			return fmt.Errorf("validation failed: %v", errs)
		}
		return fmt.Errorf("validation error: %w", err)
	}

	return nil
}

func BindWithCustomValidation(
	ctx *gin.Context,
	obj any,
	validator func(any) error,
) error {
	if err := ctx.ShouldBindJSON(obj); err != nil {
		return err
	}
	if err := validator(obj); err != nil {
		return err
	}
	return nil
}
