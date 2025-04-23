package validator

import (
	"context"
	"fmt"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/knstch/subtrack-libs/svcerrs"
	"github.com/knstch/subtrack-libs/tracing"
)

func ValidatePassword(ctx context.Context, password string) error {
	ctx, span := tracing.StartSpan(ctx, "validator: ValidatePassword")
	defer span.End()

	if err := validation.ValidateWithContext(ctx, password, validation.By(passwordValidator(password)), validation.Required); err != nil {
		return err
	}

	return nil
}

func passwordValidator(password string) validation.RuleFunc {
	return func(value interface{}) error {
		if len([]rune(password)) < 8 {
			return fmt.Errorf("password must be at least 8 characters: %w", svcerrs.ErrInvalidData)
		}

		var hasLetter, hasNumber, hasUpperCase bool

		for _, v := range password {
			switch {
			case unicode.IsLetter(v):
				hasLetter = true
				if unicode.IsUpper(v) {
					hasUpperCase = true
				}
			case unicode.IsNumber(v):
				hasNumber = true
			}
		}

		if !hasLetter || !hasNumber || !hasUpperCase {
			return fmt.Errorf("password must have letters, number, and upper case: %w", svcerrs.ErrInvalidData)
		}

		return nil
	}
}
