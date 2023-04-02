//go:build unit

package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
	"testing"
)

type toBeValidated struct {
	RequiredField       string `validate:"required"`
	GreaterThanTwoField int    `validate:"gt=2"`
}

func TestReformat(t *testing.T) {
	t.Run("return custom message for `required` tag", func(t *testing.T) {
		v := validator.New()

		err := v.Struct(toBeValidated{GreaterThanTwoField: 3})

		require.Equal(t, "RequiredField is required", Reformat(err).Error())
	})

	t.Run("return custom message for `gt` tag", func(t *testing.T) {
		v := validator.New()

		err := v.Struct(toBeValidated{
			RequiredField:       "some-value",
			GreaterThanTwoField: 1,
		})

		require.Equal(t, "GreaterThanTwoField must be greater than 2", Reformat(err).Error())
	})
}
