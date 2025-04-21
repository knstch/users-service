package validator

import (
	"fmt"
	"regexp"

	"github.com/knstch/subtrack-libs/svcerrs"
)

var (
	codeExp = regexp.MustCompile(`^\d{4}$`)
)

func ValidateConfirmationCode(code string) error {
	match := codeExp.MatchString(code)
	if !match {
		return fmt.Errorf("invalid code: %w", svcerrs.ErrInvalidData)
	}

	return nil
}
