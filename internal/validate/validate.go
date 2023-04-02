package validate

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

// Reformat provides custom messages for ValidationErrors since default errors are too verbose
func Reformat(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var formattedErrs []string

		for _, fe := range ve {
			switch fe.Tag() {
			case "required":
				formattedErrs = append(formattedErrs, fmt.Sprintf("%s is required", fe.Field()))
			case "gt":
				formattedErrs = append(formattedErrs, fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param()))
			}
		}

		return errors.New(strings.Join(formattedErrs, "\n"))
	}

	return err
}
